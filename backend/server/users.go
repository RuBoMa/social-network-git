package server

import (
	"log"
	"real-time-forum/database"
	"sort"
	"time"

	"github.com/gorilla/websocket"
)

// Broadcast the active users list exluding the user themselves
func broadcastUsers() {

	// Send sorted list to each client
	for client, userID := range clients {
		sortedUsers := sortUsers(userID)

		// Send the list of active users back to the client
		message := Message{
			Type:  "update_users",
			Users: sortedUsers, // Send the active users list
		}
		err := client.WriteJSON(message)
		if err != nil {
			log.Println("Error sending user update:", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func sendChatPartner(conn *websocket.Conn, msg Message, userID int) {

	participants, err := database.GetParticipants(msg.ChatID)
	if err != nil {
		log.Println("Issue getting participant", err)
		return
	}
	var chatPartner User
	for _, user := range participants {
		if user != userID {
			chatPartner.ID = user
		}
	}
	username, err := database.GetUsername(chatPartner.ID)
	if err != nil {
		log.Println("Error getting username", err)
		return
	}
	chatPartner.Username = username
	chatPartner.Online = false

	for _, clientID := range clients {
		if clientID == chatPartner.ID {
			// If the user ID exists in the clients map, they are online
			chatPartner.Online = true
			break
		}
	}

	message := Message{
		Type:     "user",
		ChatUser: chatPartner,
	}

	err = conn.WriteJSON(message)
	if err != nil {
		log.Println("Error sending history:", err)
	}
}

// Sorts users: latest conversations first, then alphabetically
func sortUsers(userID int) []User {
	var sortedUsers []UserInteraction
	var noInteractionUsers []User

	allUsers, err := database.GetUsers()
	if err != nil {
		log.Println("Error fetching users:", err)
		return nil
	}

	// Iterate through all clients (users)
	for user_id, username := range allUsers {

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

			// If we have a timestamp, add the user to the sorted list
			sortedUsers = append(sortedUsers, UserInteraction{
				UserID:          user_id,
				Username:        username,
				LastInteraction: interactionTime,
			})
		} else {
			noInteractionUsers = append(noInteractionUsers, User{
				ID:       user_id,
				Username: username,
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
			log.Println("Error parsing timestamp for user", sortedUsers[i].Username, errI)
			return false // or handle this case as needed
		}
		if errJ != nil {
			log.Println("Error parsing timestamp for user", sortedUsers[j].Username, errJ)
			return false // or handle this case as needed
		}

		// Compare the time objects: descending order (most recent first)
		return timestampI.After(timestampJ)
	})

	// Sort users with no interactions alphabetically
	sort.Slice(noInteractionUsers, func(i, j int) bool {
		return noInteractionUsers[i].Username < noInteractionUsers[j].Username
	})

	// Combine both lists: users with interactions first, then users without interactions
	var finalSortedUsers []User
	for _, user := range sortedUsers {
		finalSortedUsers = append(finalSortedUsers, User{
			ID:       user.UserID,
			Username: user.Username,
		})
	}
	for _, user := range noInteractionUsers {
		finalSortedUsers = append(finalSortedUsers, User{
			ID:       user.ID,
			Username: user.Username,
		})
	}
	for i, user := range finalSortedUsers {
		online := false

		// Loop through the clients map to check if this user has an active connection
		for _, clientID := range clients {
			if clientID == user.ID {
				// If the user ID exists in the clients map, they are online
				online = true
				break
			}
		}

		// Set the user's online status
		finalSortedUsers[i].Online = online
	}

	return finalSortedUsers
}

// Get the current timestamp
func GetTimestamp() int64 {
	return time.Now().Unix()
}
