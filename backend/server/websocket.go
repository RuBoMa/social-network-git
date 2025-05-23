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
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	defer func() {
		chat.CloseConnection(userID)
	}()
	log.Println("New WebSocket connection from user:", userID)

	chat.ClientsMutex.Lock()
	chat.Clients[userID] = conn
	log.Println("User added to chat clients:", userID)
	//chat.BroadcastUsers() // BROADCAST ONLY USERS WITH DISCUSSION
	chat.ClientsMutex.Unlock()

	log.Println(chat.Clients)

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
		log.Printf("Received message: %+v\n", msg)
		err = database.AddMessageIntoDB(msg.Sender.UserID, msg.Receiver.UserID, msg.GroupID, msg.Content, false)
		if err != nil {
			log.Println("Error adding message to database:", err)
			continue
		}

		message := models.ChatMessage{}

		switch msg.Type {
		case "chatBE":
			message = chat.HandleChatHistory(msg)

		case "messageBE":
			message = chat.HandleChatMessage(msg)

		case "typingBE", "stopTypingBE":
			message = chat.HandleTypingStatus(msg)
		}
		chat.Broadcast <- message
		chat.MessagesMutex.Unlock()

	}
}
