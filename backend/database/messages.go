package database

import (
	"database/sql"
	"log"
	"social_network/models"
	"time"
)

// AddMessageIntoDB inserts a new group into the database
// It takes the group name, description, creator ID, and privacy setting as parameters
func AddMessageIntoDB(senderID, receiverID, groupID int, content string, isRead bool) (int, error) {
	// log.Printf("Saving message to database: SenderID=%d, ReceiverID=%d, Content=%s\n", senderID, receiverID, content) // debug

	result, err := db.Exec("INSERT INTO Messages (sender_id, received_id, group_id, content, is_read, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		senderID, receiverID, groupID, content, isRead, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Println("Error inserting message to database:", err)
		return 0, err
	}

	messageID, err := result.LastInsertId()
	if err != nil {
		log.Println("Error getting last insert ID:", err)
		return 0, err
	}

	return int(messageID), nil
}

// GetMessageByID retrieves a message from the database by its ID
// It returns the message along with its sender and receiver information
func GetMessageByID(messageID int) (models.ChatMessage, error) {
	var message models.ChatMessage

	err := db.QueryRow("SELECT sender_id, received_id, group_id, content, is_read, created_at FROM Messages WHERE id = ?", messageID).
		Scan(&message.Sender.UserID, &message.Receiver.UserID, &message.GroupID, &message.Content, &message.IsRead, &message.CreatedAt)
	if err != nil {
		log.Println("Error retrieving message by ID:", err)
		return message, err
	}

	message.Sender, err = GetUser(message.Sender.UserID)
	if err != nil {
		log.Println("Error fetching username for id:", message.Sender.UserID)
		return message, err
	}
	message.MessageID = messageID

	return message, nil
}

// GetHistory retrieves the chat history between two users or for a group
func GetHistory(userID1, userID2, groupID int) ([]models.ChatMessage, error) {
	var rows *sql.Rows
	var err error
	chats := []models.ChatMessage{}

	if groupID == 0 {
		query := `
			SELECT sender_id, content, is_read, created_at
			FROM Messages
			WHERE group_id = 0
			  AND ((sender_id = ? AND received_id = ?) OR (sender_id = ? AND received_id = ?))
			ORDER BY created_at ASC`
		rows, err = db.Query(query, userID1, userID2, userID2, userID1)
	} else {
		query := `
			SELECT sender_id, content, is_read, created_at
			FROM Messages
			WHERE group_id = ?
			ORDER BY created_at ASC`
		rows, err = db.Query(query, groupID)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return chats, nil
		}
		return chats, err
	}
	defer rows.Close()

	for rows.Next() {
		message := models.ChatMessage{}
		if err := rows.Scan(&message.Sender.UserID, &message.Content, &message.IsRead, &message.CreatedAt); err != nil {
			return chats, err
		}
		message.Sender, err = GetUser(message.Sender.UserID)
		if err != nil {
			log.Println("Error fetching username for id: ", message.Sender.UserID)
			return chats, err
		}
		message.Sender.Email = ""       // Clear email for privacy
		message.Sender.AboutMe = ""     // Clear about me for privacy
		message.Sender.DateOfBirth = "" // Clear date of birth for privacy

		chats = append(chats, message)
	}

	if err := rows.Err(); err != nil {
		return chats, err
	}
	return chats, nil

}

// GetInterectedUsers retrieves a list of users who have interacted with a specific messageID
func GetInteractedUsersByMessageID(messageID int) ([]models.User, error) {
	var rows *sql.Rows
	var err error
	interactedUsers := []models.User{}

	query := `
		SELECT sender_id, received_id, group_id
		FROM Messages
		WHERE id = ?`
	rows, err = db.Query(query, messageID)
	if err != nil {
		log.Println("Error fetching interacted users by message ID:", err)
		return interactedUsers, err
	}
	defer rows.Close()

	for rows.Next() {
		message := models.ChatMessage{}

		if err := rows.Scan(&message.Sender.UserID, &message.Receiver.UserID, &message.GroupID); err != nil {
			return interactedUsers, err
		}
		if message.GroupID != 0 {
			interactedUsers, err = GetGroupMembers(message.GroupID)
			if err != nil {
				log.Println("Error fetching group members for group ID:", message.GroupID)
				return interactedUsers, err
			}
		} else {
			interactedUsers = append(interactedUsers, message.Sender)
			interactedUsers = append(interactedUsers, message.Receiver)
		}
	}

	if err := rows.Err(); err != nil {
		return interactedUsers, err
	}
	return interactedUsers, nil
}

// GetMessage retrieves a message from the database by its ID
// func GetMessage(message_id int) ([]string, error) {
// 	var message []string
// 	var chatID int
// 	var senderID int
// 	var content string
// 	var createdAt string

// 	err := db.QueryRow("SELECT chat_id, sender_id, content, created_at FROM Messages WHERE id = ?", message_id).Scan(&chatID, &senderID, &content, &createdAt)
// 	if err != nil {
// 		return message, err
// 	}

// 	username, err := GetUsername(senderID)
// 	if err != nil {
// 		log.Println("Error fetching username for id: ", senderID)
// 		return message, err
// 	}

// 	message = []string{
// 		fmt.Sprint(chatID),
// 		fmt.Sprint(senderID),
// 		username,
// 		content,
// 		createdAt,
// 	}

// 	return message, nil
// }

// GroupChatExists checks if a group chat exists in the database
func GroupChatExists(groupID int) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM Messages WHERE group_id = ?)", groupID).Scan(&exists)
	if err != nil {
		log.Println("Error checking if group chat exists:", err)
		return false, err
	}
	return exists, nil
}

// HasExistingConversation checks if a conversation exists between two users
func HasExistingConversation(userID1, userID2 int) (bool, error) {
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM Messages
			WHERE (
			(sender_id = ? AND received_id = ?) 
			OR (sender_id = ? AND received_id = ?)
			) AND group_id = 0
		)`, userID1, userID2, userID2, userID1).Scan(&exists)
	if err != nil {
		log.Println("Error checking existing conversation:", err)
		return false, err
	}
	return exists, nil
}
