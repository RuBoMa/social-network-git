package chat

import (
	"log"
	"social_network/database"
	"social_network/models"
	"sync"
)

var (
	Clients       = make(map[int]*models.Client)  // Map of WebSocket userID -> connection
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
		for id, client := range Clients {
			// Validate connection
			if client == nil {
				log.Printf("Connection for user %d is nil. Removing from Clients map.", id)
				delete(Clients, id)
				continue
			}

			for _, receiver := range receivers {
				if id == receiver.UserID {
					client.Mu.Lock()
					err := client.Conn.WriteJSON(message)
					client.Mu.Unlock()
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
	for userID, client := range Clients {
		log.Println("Sending active users to client:", userID)
		sortedUsers := SortUsers(userID)

		// Send the list of active users back to the client
		message := models.ChatMessage{
			Type:  "interacted_users_response",
			Users: sortedUsers, // Send the active users list
		}
		err := client.Conn.WriteJSON(message)
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

	for userID, client := range Clients {
		if notification.UserID == userID {
			err := client.Conn.WriteJSON(notification)
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
	if client, ok := Clients[userID]; ok {
		client.Conn.Close()
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
		response.Type = "typing"
	} else {
		response.Type = "stop_typing"
	}

	if msg.GroupID != 0 {
		// Handle typing status for group chats
		groupMembers, err := database.GetGroupMembers(msg.GroupID)
		if err != nil {
			log.Println("Error fetching group members:", err)
			return response
		}

		// Fetch all users
		users, err := database.GetUsers()
		if err != nil {
			log.Println("Error fetching users:", err)
			return response
		}
		var senderUser models.User
		for _, user := range users {
			if user.UserID == msg.Sender.UserID {
				senderUser = user // Get the sender's user info
				break
			}
		}

		// Pass all user info directly
		for _, member := range groupMembers {
			if member.UserID != msg.Sender.UserID { // Exclude the sender
				response.Sender = senderUser // Include sender details from the message
				response.Users = users       // Include all user info
				response.Receiver = member
				Broadcast <- response
			}
		}
	} else {
		// Handle typing status for private chats
		users, err := database.GetUsers()
		if err != nil {
			log.Println("Error fetching users:", err)
			return response
		}

		// Pass all user info directly
		response.Sender = msg.Sender // Include sender details from the message
		response.Users = users       // Include all user info
		response.Receiver = msg.Receiver
		Broadcast <- response
	}

	return response
}
