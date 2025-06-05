package chat

import (
	"log"
	"social_network/database"
	"social_network/models"
)

// HandleChatHistory retrieves the chat history between two users or for a group
func HandleChatHistory(msg models.ChatMessage) models.ChatMessage {
	chatMessage := models.ChatMessage{}

	history, err := database.GetHistory(msg.Sender.UserID, msg.Receiver.UserID, msg.GroupID)
	if err != nil {
		log.Println("Error retreiving chat history: ", err)
		return chatMessage
	}
	var isAllowed bool

	if msg.GroupID == 0 {
		isAllowed, err = database.IsEitherOneFollowing(msg.Sender.UserID, msg.Receiver.UserID)
		if err != nil {
			log.Println("Sender is not allowed to view the chat history:", msg.Sender.UserID, msg.Receiver.UserID)
			chatMessage.Type = "error"
			chatMessage.Content = "Cannot check follow status"
			return chatMessage
		}
	} else {
		isAllowed = true
	}

	chatMessage = models.ChatMessage{
		Type:    "chat",
		History: history,
		Receiver: models.User{
			UserID: msg.Receiver.UserID,
		},
		GroupID:    msg.GroupID,
		CanMessage: isAllowed,
	}

	return chatMessage
}

// HandleChatMessage adds the message to the database and return is with the type "message"
func HandleChatMessage(msg models.ChatMessage) models.ChatMessage {

	// Check from database if either one if following the other
	message := msg
	if msg.Sender.UserID == 0 || (msg.Receiver.UserID == 0 && msg.GroupID == 0) {
		log.Println("Invalid sender, receiver, or group:", msg)
		message.Type = "error"
		message.Content = "Invalid sender, receiver, or group"
		return message
	}

	// If it's a group message, ensure the group exists and the sender is part of the group
	if msg.GroupID != 0 {
		isMember, err := database.IsGroupMember(msg.Sender.UserID, msg.GroupID)
		if err != nil || !isMember {
			log.Println("Sender is not a member of the group:", msg.GroupID)
			message.Type = "error"
			message.Content = "Sender is not a member of the group"
			return message
		}
	} else {
		isAllowed, err := database.IsEitherOneFollowing(msg.Sender.UserID, msg.Receiver.UserID)
		if err != nil || !isAllowed {
			log.Println("Sender is not allowed to message the receiver:", msg.Sender.UserID, msg.Receiver.UserID)
			message.Type = "error"
			message.Content = "One of the users need to follow the other to send messages"
			return message
		}
	}

	// Add the message to the database
	messageID, err := database.AddMessageIntoDB(msg.Sender.UserID, msg.Receiver.UserID, msg.GroupID, msg.Content, false)
	if err != nil {
		log.Println("Message not added to database:", err)
		message.Type = "error"
		message.Content = "Failed to save message"
		return message
	}

	messageDB, err := database.GetMessageByID(messageID)
	if err != nil {
		log.Println("Error fetching message from database:", err)
		message.Type = "error"
		message.Content = "Failed to fetch saved message"
		return message
	}

	message.Sender = messageDB.Sender
	message.CreatedAt = messageDB.CreatedAt
	message.MessageID = messageDB.MessageID

	message.Type = "message"
	return message
}
