package chat

import (
	"log"
	"social_network/database"
	"social_network/models"
	"sync"
)

var (
	Clients       = make(map[int]*models.Client)  // Map of WebSocket userID -> connection
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
				if message.Type != "typing" && message.Type != "stopTyping" {
					receivers = append(receivers, message.Sender)
				}
			}
		}

		for id, client := range Clients {
			// Validate connection
			if client == nil {
				log.Printf("Connection for user %d is nil. Removing from Clients map.", id)
				CloseConnection(id)
				continue
			}

			for _, receiver := range receivers {
				if id == receiver.UserID {
					//client.Mu.Lock()
					ClientsMutex.Lock()

					err := client.Conn.WriteJSON(message)
					ClientsMutex.Unlock()

					//client.Mu.Unlock()
					if err != nil {
						log.Printf("Error writing message to user %d: %v. Closing connection.", id, err)
						CloseConnection(id)
					}
				}
			}
		}
	}
}

// BroadcastNotification sends a notification to the defined client if they are online
// Front can listen notifications based on the type of notification
func BroadcastNotification(notification models.Notification) {

	for userID, client := range Clients {
		if notification.UserID == userID {
			ClientsMutex.Lock()
			err := client.Conn.WriteJSON(notification)
			ClientsMutex.Unlock()
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

	ClientsMutex.Lock()
	if client, ok := Clients[userID]; ok {
		client.Conn.Close()
		delete(Clients, userID)
	}
	ClientsMutex.Unlock()
}

// HandleTypingStatus handles the typing status of users in a chat
// It sends a typing or stop typing message to the receiver
func HandleTypingStatus(msg models.ChatMessage) {
	response := models.ChatMessage{}
	var err error

	if msg.Type == "typingBE" {
		response.Type = "typing"

	} else {
		response.Type = "stop_typing"
	}

	response.Sender, err = database.GetUser(msg.Sender.UserID)
	if err != nil {
		log.Println("Error fetching sender user details:", err)
	}

	if msg.GroupID != 0 {
		response.GroupID = msg.GroupID
		// Handle typing status for group chats
		groupMembers, err := database.GetGroupMembers(msg.GroupID)
		if err != nil {
			log.Println("Error fetching group members:", err)
		}

		// Pass all user info directly
		for _, member := range groupMembers {
			if member.UserID == msg.Sender.UserID {
				continue // Exclude the sender
			}
			response.Receiver = member
			Broadcast <- response

		}
	} else {

		response.Receiver = msg.Receiver
		Broadcast <- response
	}
}
