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
	chatMessage := models.ChatMessage{}

	history, err := database.GetHistory(msg.Sender.UserID, msg.Receiver.UserID, msg.GroupID)
	if err != nil {
		log.Println("Error retreiving chat history: ", err)
		return chatMessage
	}

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
	return message

}

// Sorts users: latest conversations first, then alphabetically
func SortUsers(userID int) []models.User {
	var sortedUsers []models.UserInteraction
	var noInteractionUsers []models.User

	allUsers, err := database.GetUsers()
	if err != nil {
		log.Println("Error fetching users:", err)
		return nil
	}

	// Iterate through all clients (users)
	for _, user := range allUsers {
		user_id := user.UserID
		username := user.Nickname

		// Skip the current user
		if user_id == userID {
			continue
		}

		// // Check for interactions where the current user is involved (either as the user or as the other user)
		interactionTime, err := database.GetLastAction(userID, user_id)
		if err != nil {
			log.Println("Error fetching latest activity:", err)
			return nil
		}
		if interactionTime != "" {
			user := models.User{
				UserID:   user_id,
				Nickname: username,
			}

			// If we have a timestamp, add the user to the sorted list
			sortedUsers = append(sortedUsers, models.UserInteraction{
				User:            user,
				LastInteraction: interactionTime,
			})
		} else {
			noInteractionUsers = append(noInteractionUsers, models.User{
				UserID:   user_id,
				Nickname: username,
			})
		}
	}

	// Sort users with interactions by the last interaction timestamp (descending)
	sort.Slice(sortedUsers, func(i, j int) bool {
		layout := "2006-01-02 15:04:05" // The format you're using

		timestampI, errI := time.Parse(layout, sortedUsers[i].LastInteraction)
		timestampJ, errJ := time.Parse(layout, sortedUsers[j].LastInteraction)

		// Handle parsing errors (optional, depending on your needs)
		if errI != nil {
			log.Println("Error parsing timestamp for user", sortedUsers[i].User.Nickname, errI)
			return false // or handle this case as needed
		}
		if errJ != nil {
			log.Println("Error parsing timestamp for user", sortedUsers[j].User.Nickname, errJ)
			return false // or handle this case as needed
		}

		// Compare the time objects: descending order (most recent first)
		return timestampI.After(timestampJ)
	})

	// Sort users with no interactions alphabetically
	sort.Slice(noInteractionUsers, func(i, j int) bool {
		return noInteractionUsers[i].Nickname < noInteractionUsers[j].Nickname
	})

	// Combine both lists: users with interactions first, then users without interactions
	var finalSortedUsers []models.User
	for _, user := range sortedUsers {
		finalSortedUsers = append(finalSortedUsers, models.User{
			UserID:   user.User.UserID,
			Nickname: user.User.Nickname,
		})
	}
	for _, user := range noInteractionUsers {
		finalSortedUsers = append(finalSortedUsers, models.User{
			UserID:   user.UserID,
			Nickname: user.Nickname,
		})
	}

	return finalSortedUsers
}

// Get the current timestamp
func GetTimestamp() int64 {
	return time.Now().Unix()
}
