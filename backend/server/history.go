package server

import (
	"log"
	"real-time-forum/database"
	"github.com/gorilla/websocket"
)

func HandleChatHistory(conn *websocket.Conn, userID int, msg Message) {
	chatUser := msg.ChatUser
	chatUser.Online = false

	// Loop through the clients map to check if this user has an active connection
	for _, clientID := range clients {
		if clientID == chatUser.ID {
			// If the user ID exists in the clients map, they are online
			chatUser.Online = true
			break
		}
	}

	chatID, err := database.GetChatID(userID, chatUser.ID)
	if err != nil {
		log.Println("Error getting chatID: ", err)
		return
	}

	var history []map[string]interface{}
	err = database.GetHistory(chatID, &history)
	if err != nil {
		log.Println("Error retreiving chat history: ", err)
		return
	}

	// Convert []map[string]interface{} to []Message
	var messages []Message
	for _, entry := range history {
		sender, _ := entry["senderID"].(int) // Convert sender to int
		username, _ := entry["senderUsername"].(string)
		content, _ := entry["content"].(string)     // Convert content to string
		createdAt, _ := entry["createdAt"].(string) // Convert timestamp to string

		messages = append(messages, Message{
			Sender: User{
				ID:       sender,
				Username: username,
			},
			Content:   content,
			CreatedAt: createdAt,
		})
	}

	message := Message{
		Type:     "chat",
		History:  messages,
		ChatID:   chatID,
		ChatUser: chatUser,
	}
	err = conn.WriteJSON(message)
	if err != nil {
		log.Println("Error sending history:", err)
	}
}
