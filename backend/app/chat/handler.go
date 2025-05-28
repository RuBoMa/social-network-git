package chat

import (
	"log"
	"social_network/database"
	"social_network/models"
	"sort"
	"time"
)

// HandleChatHistory retrieves the chat history between two users or for a group
func HandleChatHistory(msg models.ChatMessage) models.ChatMessage {
	log.Println("Retrieving chat history for users:", msg.Sender.UserID, "and", msg.Receiver.UserID, "in group:", msg.GroupID)
	chatMessage := models.ChatMessage{}

	history, err := database.GetHistory(msg.Sender.UserID, msg.Receiver.UserID, msg.GroupID)
	if err != nil {
		log.Println("Error retreiving chat history: ", err)
		return chatMessage
	}

	log.Printf("Chat history retrieved: %+v\n", history)

	chatMessage = models.ChatMessage{
		Type:    "chat",
		History: history,
		Receiver: models.User{
			UserID: msg.Receiver.UserID,
		},
	}

	return chatMessage
}

// HandleChatMessage adds the message to the database and return is with the type "message"
func HandleChatMessage(msg models.ChatMessage) models.ChatMessage {
	log.Println("are we getting here")
	message := msg
	if msg.Sender.UserID == 0 || msg.Receiver.UserID == 0 {
		log.Println("Invalid sender or receiver:", msg)
		message.Type = "error"
		message.Content = "Invalid sender or receiver"
		return message
	}
	// Add the message to the database
	err := database.AddMessageIntoDB(msg.Sender.UserID, msg.Receiver.UserID, msg.GroupID, msg.Content, false)
	if err != nil {
		log.Println("Message not added to database:", err)
		return message
	}
	message.Type = "message"
	log.Println("Message successfully saved to database:", message)

	go BroadcastSortedUsers(msg.Sender.UserID)
	go BroadcastSortedUsers(msg.Receiver.UserID)

	return message
}

// Sorts users: latest conversations first, then alphabetically
func SortUsers(userID int) []models.User {
	// 1. Get all users the current user has chatted with, with last interaction timestamp
	interactions, err := database.GetUserInteractions(userID)
	if err != nil {
		return []models.User{}
	}

	// 2. Sort by last interaction (descending)
	sort.SliceStable(interactions, func(i, j int) bool {
		return interactions[i].LastInteraction > interactions[j].LastInteraction
	})

	var sortedUsers []models.User

	// 3. For each interaction, fetch full user info
	for _, interaction := range interactions {
		otherUserID := interaction.User.UserID
		if otherUserID == userID {
			continue // skip self
		}
		fullUser, err := database.GetUser(otherUserID)
		if err != nil {
			continue // skip if can't fetch
		}
		sortedUsers = append(sortedUsers, fullUser)
	}

	// 4. Optionally, sort alphabetically by nickname for users with same last interaction
	sort.SliceStable(sortedUsers, func(i, j int) bool {
		if interactions[i].LastInteraction == interactions[j].LastInteraction {
			return sortedUsers[i].Nickname < sortedUsers[j].Nickname
		}
		return false
	})
	log.Printf("Sorted users for user %d: %+v\n", userID, sortedUsers)
	
	return sortedUsers
}

// Get the current timestamp
func GetTimestamp() int64 {
	return time.Now().Unix()
}

func HandleChatInitiation(msg models.ChatMessage) models.ChatMessage {
	log.Printf("Initiating chat between users: %d and %d\n", msg.Sender.UserID, msg.Receiver.UserID)

	// Check if users can chat (following relationship)
	canChat, err := database.CanUsersChat(msg.Sender.UserID, msg.Receiver.UserID)
	if err != nil || !canChat {
		log.Println("Users cannot chat - no following relationship")
		return models.ChatMessage{
			Type:    "error",
			Content: "Cannot start chat - users must follow each other",
		}
	}

	// Add an initial "chat initiated" message to create the chat history entry
	err = database.AddMessageIntoDB(msg.Sender.UserID, msg.Receiver.UserID, 0, "Chat initiated", true)
	if err != nil {
		log.Println("Error creating initial chat message:", err)
		return models.ChatMessage{
			Type:    "error",
			Content: "Failed to initiate chat",
		}
	}
	go BroadcastSortedUsers(msg.Sender.UserID)
	go BroadcastSortedUsers(msg.Receiver.UserID)

	// Return success message
	return models.ChatMessage{
		Type:     "chat_initiated",
		Sender:   msg.Sender,
		Receiver: msg.Receiver,
		Content:  "Chat initiated successfully",
	}
}
