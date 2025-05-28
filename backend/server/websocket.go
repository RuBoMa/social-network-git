package server

import (
	"encoding/json"
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
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// upgrade to Websocket protocol
	log.Println("Attempting to upgrade connection to WebSocket")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	log.Println("WebSocket connection upgraded successfully")
	defer func() {
		chat.CloseConnection(userID)
	}()
	log.Println("New WebSocket connection from user:", userID)

	chat.ClientsMutex.Lock()
	chat.Clients[userID] = conn
	log.Printf("User %d added to chat clients. Current clients: %+v\n", userID, chat.Clients)
	log.Println("User added to chat clients:", userID)
	// chat.BroadcastUsers() // BROADCAST ONLY USERS WITH DISCUSSION, change broadcast logic
	chat.ClientsMutex.Unlock()

	go chat.BroadcastSortedUsers(userID) // Send sorted users to the new client
	log.Println("we are here:", chat.Clients)

	var msg models.ChatMessage

	// Indefinite loop to listen messages while connection open
	for {
		log.Println("Waiting for message...")
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			chat.CloseConnection(userID)
			break
		}
		log.Println("Message received:", string(p))
		chat.MessagesMutex.Lock()

		err = json.Unmarshal(p, &msg) // Unmarshal the bytes into the struct
		if err != nil {
			log.Println("Error unmarshalling JSON:", err)
			continue // Currently not crashing the server, invalid message format will be ignored
		}
		msg.Sender.UserID = userID
		log.Printf("Received message: %+v\n", msg)

		message := models.ChatMessage{}

		switch msg.Type {
		case "chatBE":
			message = chat.HandleChatHistory(msg)
			// conn.WriteJSON(historyMsg) // Send chat history back to the client

		case "message":
			log.Println("Handling message")
			message := chat.HandleChatMessage(msg)
			go chat.BroadcastSortedUsers(userID)
			go chat.BroadcastSortedUsers(msg.Receiver.UserID)
			chat.Broadcast <- message

		case "typingBE", "stopTypingBE":
			message := chat.HandleTypingStatus(msg)
			chat.Broadcast <- message

		case "mark_notification_read":
			log.Println("Marking notification as read:", msg.NotificationID)
			// Mark notification as read
			if msg.NotificationID != 0 {
				err := database.NotificationSeen(msg.NotificationID)
				if err == nil {
					chat.Broadcast <- msg
				}
			}

			message = chat.HandleTypingStatus(msg)
			// chat.Broadcast <- message
		case "requestSortedUsers":
			go chat.BroadcastSortedUsers(userID) // Send sorted users to the client
			continue                             // No need to broadcast, just send sorted users
		case "initiate_chat":
			log.Println("Handling chat initiation")
			message := chat.HandleChatInitiation(msg)
			if message.Type != "error" {
				go chat.BroadcastSortedUsers(userID)
				go chat.BroadcastSortedUsers(msg.Receiver.UserID)
			}
			chat.Broadcast <- message
		}
		chat.Broadcast <- message // Broadcast the message to all clients
		chat.MessagesMutex.Unlock()
	}
}
