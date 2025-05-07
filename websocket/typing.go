package websocket

import (
	"real-time-forum/backend"
)

func sendTypingStatus(msg Message, userID int) {

	response := Message{
		Type:   "typing",
		ChatID: msg.ChatID,
	}

	if msg.Type == "stopTypingBE" {
		response.Type = "stop_typing"
	}

	chatParties, err := backend.GetParticipants(msg.ChatID)
	if err != nil {
		return
	}

	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for i, clientID := range clients {
		for _, chatUser := range chatParties {
			// If the user ID exists in the clients map, they are online
			if clientID == userID {
				continue
			}
			if clientID == chatUser {
				err := i.WriteJSON(response)
				if err != nil {
					return
				}
			}
		}

	}
}
