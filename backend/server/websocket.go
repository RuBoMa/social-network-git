package server

import (
	"encoding/json"
	"log"
	"net/http"
	"social_network/app"
	"social_network/app/chat"
	"social_network/models"

	"github.com/gorilla/websocket"
)

// updates HTTP connection to websocket protocol
// checks if the request is from localhost:3000 (frontend)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
//origin := r.Header.Get("Origin")
//origin == "http://localhost:3000"

// Handles Websocket connections
func HandleConnections(w http.ResponseWriter, r *http.Request) {

	loggedIn, userID := app.VerifySession(r)
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
	defer conn.Close()
	log.Println("New websocket connection established")

	// defer func() {
	// 	// Close the connection properly
	// 	chat.CloseConnection(userID)
	// }()

	chat.ClientsMutex.Lock()
	// add user to clients
	chat.Clients[userID] = conn
	chat.BroadcastUsers() // DISCUSS LOGIC WITH THE GROUP
	chat.ClientsMutex.Unlock()

	var msg models.ChatMessage

	// Indefinite loop to listen messages while connection open
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			chat.CloseConnection(userID)
			break
		}

		chat.MessagesMutex.Lock()

		err = json.Unmarshal(p, &msg) // Unmarshal the bytes into the struct
		if err != nil {
			log.Println("Error unmarshalling JSON:", err)
			continue // Currently not crashing the server, invalid message format will be ignored
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
