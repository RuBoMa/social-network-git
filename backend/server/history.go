package server

import (
	"log"
	"social_network/database"

	"github.com/gorilla/websocket"
)

func HandleChatHistory(conn *websocket.Conn, msg Message) {
	senderID := msg.Sender.ID
	receiverID := msg.Receiver.ID

	var history []map[string]interface{}
	err := database.GetHistory(senderID, receiverID, 0, &history)
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
		isRead, _ := entry["isRead"].(bool)         // Convert isRead to bool
		createdAt, _ := entry["createdAt"].(string) // Convert timestamp to string

		messages = append(messages, Message{
			Sender: User{
				ID:       sender,
				Username: username,
			},
			Content:   content,
			IsRead:    isRead,
			CreatedAt: createdAt,
		})
	}

	message := Message{
		Type:    "chat",
		History: messages,
		Receiver: User{
			ID: receiverID,
		},
	}
	err = conn.WriteJSON(message)
	if err != nil {
		log.Println("Error sending history:", err)
	}
}
