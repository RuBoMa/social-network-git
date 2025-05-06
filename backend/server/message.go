package server

import (
	"log"
	"real-time-forum/backend/database"
)

// Broadcast messages to all active clients
func BroadcastMessages() {
	for {
		message := <-broadcast
		participants, err := database.GetParticipants(message.ChatID)
		if err != nil {
			log.Println(err)
			return
		}
		clientsMutex.Lock()
		for client, id := range clients {
			for _, user := range participants {
				if user == id {
					err := client.WriteJSON(message)
					if err != nil {
						log.Println("Write error:", err)
						client.Close()
						delete(clients, client)
					}
				}
			}
		}
		clientsMutex.Unlock()
		broadcastUsers()
	}
}

func AddChatToDB(userID int, msg Message) int {

	message_id, err := database.AddMessageToDB(userID, msg.Content, msg.ChatID)
	if err != nil {
		log.Println("Error adding message:", err)
		return 0
	}

	return message_id
}
