package chat

import (
	"log"
	"social_network/database"
	"social_network/models"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	Clients       = make(map[int]*websocket.Conn) // Map of WebSocket userID -> connection
	typingMap     = make(map[int]bool)            // Map of userID -> typing status
	Broadcast     = make(chan models.ChatMessage) // Channel for broadcasting messages
	ClientsMutex  sync.Mutex                      // Protects access to activeUsers map
	MessagesMutex sync.Mutex                      // Protects access to broadcast channel
)

// Broadcast messages to relevant users
// This function runs in a separate goroutine and listens for messages on the broadcast channel.
func BroadcastMessages() {
	log.Println("Starting message broadcast loop...")
	for {
		message := <-Broadcast
		var receivers []models.User
		var err error

		if message.GroupID != 0 {
			// Get all group members from the database
			receivers, err = database.GetGroupMembers(message.GroupID)
			if err != nil {
				log.Println("Error fetching group members:", err)
				continue
			}
		} else {
			// Add sender and receiver for private messages
			if message.Receiver.UserID != 0 {
				receivers = append(receivers, message.Receiver)
			}
			if message.Sender.UserID != 0 {
				receivers = append(receivers, message.Sender)
			}
		}

		log.Printf("Receivers for message: %+v\n", receivers)

		ClientsMutex.Lock()
		for id, conn := range Clients {
			// Validate connection
			if conn == nil {
				log.Printf("Connection for user %d is nil. Removing from Clients map.", id)
				delete(Clients, id)
				continue
			}

			for _, receiver := range receivers {
				if id == receiver.UserID {
					err := conn.WriteJSON(message)
					if err != nil {
						log.Printf("Error writing message to user %d: %v. Closing connection.", id, err)
						CloseConnection(id)
					}
				}
			}
		}
		ClientsMutex.Unlock()
		// BroadcastUsers()
	}
}
func GetInteractedUsers(userID int) ([]models.User, error) {
	return database.GetInteractedUsers(userID)
}

// Broadcast the active users list exluding the user themselves
// DISCUSS LOGIC WITH THE GROUP
func BroadcastUsers() {
	log.Println("Broadcasting active users to clients...")
	ClientsMutex.Lock()
	defer ClientsMutex.Unlock()

	// Send sorted list to each client
	for userID, conn := range Clients {
		log.Println("Sending active users to client:", userID)
		sortedUsers := SortUsers(userID)

		// Send the list of active users back to the client
		message := models.ChatMessage{
			Type:  "interacted_users",
			Users: sortedUsers, // Send the active users list
		}
		err := conn.WriteJSON(message)
		if err != nil {
			log.Println("Error sending user update:", err)
			CloseConnection(userID)
		}
	}
	log.Printf("Current clients: %+v\n", Clients)
}

// BrioadcastNotification sends a notification to the defined clien if they are online
// Front can listen notifications based on the type of notification
func BroadcastNotification(notification models.Notification) {
	log.Println("Broadcasting notification to clients...")
	ClientsMutex.Lock()
	defer ClientsMutex.Unlock()

	for userID, conn := range Clients {
		if notification.UserID == userID {
			err := conn.WriteJSON(notification)
			if err != nil {
				log.Println("Error sending notification:", err)
				CloseConnection(userID)
			}
		}
	}
}

// CloseConnection closes the WebSocket connection properly for a user
func CloseConnection(userID int) {
	log.Println("Closing connection for user:", userID)
	// Send stop typing status to all user who are getting typing status
	for id, isTyping := range typingMap {
		if isTyping {
			message := HandleTypingStatus(models.ChatMessage{
				Type: "stopTypingBE",
				Receiver: models.User{
					UserID: id,
				},
			})
			Broadcast <- message
		}
	}

	ClientsMutex.Lock()
	if conn, ok := Clients[userID]; ok {
		conn.Close()
		delete(Clients, userID)
	}
	ClientsMutex.Unlock()
	// BroadcastUsers()
}

// HandleTypingStatus handles the typing status of users in a chat
// It sends a typing or stop typing message to the receiver
func HandleTypingStatus(msg models.ChatMessage) models.ChatMessage {

	response := models.ChatMessage{}

	if msg.Type == "typingBE" {
		if typingMap[msg.Receiver.UserID] {
			return response
		}
		response.Type = "typing"
	} else {
		response.Type = "stop_typing"
	}

	// Toggle the typing status
	typingMap[msg.Receiver.UserID] = !typingMap[msg.Receiver.UserID]
	response.Receiver = msg.Receiver

	return response

}