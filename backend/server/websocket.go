package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"social_network/app"
	"social_network/app/chat"
	"social_network/database"
	"social_network/models"

	"github.com/gorilla/websocket"
)

// updates HTTP connection to websocket protocol
// checks if the request is from localhost:3000 (frontend)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "http://localhost:3000"
	},
}

// Handles Websocket connections
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	loggedIn, userID := app.VerifySessionToken(token)
	if !loggedIn {
		app.ResponseHandler(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// upgrade to Websocket protocol
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer func() {
		chat.CloseConnection(userID)
	}()

	chat.ClientsMutex.Lock()
	chat.Clients[userID] = &models.Client{
		Conn: conn,
	}
	chat.ClientsMutex.Unlock()
	SendInteractedUsers(userID) // Send interacted users to the new connection

	var msg models.ChatMessage

	// Indefinite loop to listen messages while connection open
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			chat.CloseConnection(userID)
			break
		}
		chat.MessagesMutex.Lock()

		err = json.Unmarshal(p, &msg) // Unmarshal the bytes into the struct
		if err != nil {
			log.Println("Error unmarshalling JSON:", err)
			continue // Currently not crashing the server, invalid message format will be ignored
		}
		msg.Sender.UserID = userID

		switch msg.Type {

		case "chat":
			log.Println("Handling chat history request:", msg)
			historyMsg := chat.HandleChatHistory(msg)
			conn.WriteJSON(historyMsg) // Send chat history back to the client

		case "message":
			message := chat.HandleChatMessage(msg)
			chat.Broadcast <- message
			SendInteractedUsers(userID) // Update interacted users after sending a message

		case "typingBE", "stopTypingBE":
			message := chat.HandleTypingStatus(msg)
			err := sendToUser(msg.Receiver.UserID, message)
			if err != nil {
				log.Println("Error sending typing status to user:", err)
			}
		case "ping":
			chat.MessagesMutex.Unlock()
			continue
		case "mark_notification_read":
			// Mark notification as read
			if msg.NotificationID != 0 {
				err := database.NotificationSeen(msg.NotificationID)
				if err == nil {
					chat.Broadcast <- msg
				}
			}
		}

		chat.MessagesMutex.Unlock()
	}
}

// SendInteractedUsers retrieves and sends the list of interacted users and groups to the client
func SendInteractedUsers(userID int) {
	interactedUsers, err := database.GetInteractedUsers(userID)
	if err != nil {
		log.Println("Error fetching interacted users:", err)
		return
	}

	interactedGroups, err := database.GetInteractedGroups(userID)
	if err != nil {
		log.Println("Error fetching interacted groups:", err)
		return
	}

	chat.ClientsMutex.Lock()
	defer chat.ClientsMutex.Unlock()
	client, ok := chat.Clients[userID]

	if !ok {
		log.Printf("Client not found for user %d", userID)
		return
	}

	err = client.Conn.WriteJSON(models.ChatMessage{
		Type:   "interacted_users_response",
		Users:  interactedUsers,
		Groups: interactedGroups,
	})
	if err != nil {
		log.Println("Error sending interacted users list:", err)
		chat.CloseConnection(userID)
		return
	}
}

// sendToUser sends a message to a specific user via WebSocket
func sendToUser(userID int, message interface{}) error {
	chat.ClientsMutex.Lock()
	defer chat.ClientsMutex.Unlock()
	client, ok := chat.Clients[userID]

	if !ok {
		return fmt.Errorf("no WebSocket connection found for user %d", userID)
	}

	return client.Conn.WriteJSON(message)
}
