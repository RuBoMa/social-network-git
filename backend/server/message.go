package server

import (
	"log"
	"social_network/database"
)

// Broadcast messages to all active clients
func BroadcastMessages() {
	for {
		message := <-broadcast
		clientsMutex.Lock()
		for client, id := range clients {
			if message.Sender.ID == id || message.Receiver.ID == id {
				err := client.WriteJSON(message)
				if err != nil {
					log.Println("Write error:", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
		clientsMutex.Unlock()
		broadcastUsers()
	}
}

func AddChatToDB(msg Message) int {

	message_id, err := database.AddMessageIntoDB(msg.Sender.ID, msg.Receiver.ID, 0, msg.Content, false)
	if err != nil {
		log.Println("Error adding message:", err)
		return 0
	}

	return message_id
}
