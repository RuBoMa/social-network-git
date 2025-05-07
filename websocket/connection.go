package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"social_network/backend"
	"strconv"
)

// Handles Websocket connections
func HandleConnections(w http.ResponseWriter, r *http.Request) {

	loggedIn, userID := backend.VerifySession(r)
	if !loggedIn {
		return
	}

	// upgrade to Websocket protocol
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer func() {
		// Remove the connection from the clients
		delete(clients, conn)
		conn.Close()
	}()

	clientsMutex.Lock()
	// add user to clients
	clients[conn] = userID
	broadcastUsers()
	clientsMutex.Unlock()

	var msg Message

	// Indefinite loop to listen messages while connection open
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {

			// Remove user connection & remove from active users
			clientsMutex.Lock()
			delete(clients, conn)
			broadcastUsers()
			clientsMutex.Unlock()
			break
		}

		messagesMutex.Lock()

		err = json.Unmarshal(p, &msg) // Unmarshal the bytes into the struct
		if err != nil {
			log.Println("Error unmarshalling JSON:", err)
			continue // Handle the error appropriately
		}

		switch msg.Type {
		case "chatBE":
			HandleChatHistory(conn, userID, msg)

		case "messageBE":
			messageID := AddChatToDB(userID, msg)
			if messageID != 0 {
				latestMessage, err := backend.GetMessage(messageID)
				if err != nil {
					log.Println("Error getting latest message:", err)
					return
				}
				chatID, _ := strconv.Atoi(latestMessage[0])
				senderID, _ := strconv.Atoi(latestMessage[1])
				response := Message{
					Type:   "stopTypingBE",
					ChatID: chatID,
					ChatUser: User{
						ID:       senderID,
						Username: latestMessage[2],
						Online:   true,
					},
				}
				sendTypingStatus(response, userID)
				message := Message{
					Type:   "message",
					ChatID: chatID,
					Sender: User{
						ID:       senderID,
						Username: latestMessage[2],
					},
					Content:   latestMessage[3],
					CreatedAt: latestMessage[4],
				}
				broadcast <- message
			}
		case "userBE":
			sendChatPartner(conn, msg, userID)
		case "typingBE", "stopTypingBE":
			sendTypingStatus(msg, userID)
		}
		msg = Message{}
		messagesMutex.Unlock()

	}
}
