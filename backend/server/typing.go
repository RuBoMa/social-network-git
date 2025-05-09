package server

func sendTypingStatus(msg Message) {

	response := Message{
		Type: "typing",
	}

	if msg.Type == "stopTypingBE" {
		response.Type = "stop_typing"
	}
	userID := msg.Sender.ID
	receiver := msg.Receiver.ID
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for i, clientID := range clients {

		if clientID == userID {
			continue
		}
		if clientID == receiver {
			err := i.WriteJSON(response)
			if err != nil {
				return
			}
		}

	}
}
