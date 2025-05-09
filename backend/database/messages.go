package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// AddGroupIntoDB inserts a new group into the database
// It takes the group name, description, creator ID, and privacy setting as parameters
func AddMessageIntoDB(senderID, receiverID, groupID int, content string, isRead bool) (int, error) {

	var result sql.Result
	var err error
	result, err = db.Exec("INSERT INTO Message (sender_id, receiver_id, group_id, content, is_read, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		senderID, receiverID, groupID, content, isRead, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Println("Error inserting post:", err)
		return 0, err
	}

	// Get the post id for the post inserted
	msgID, err := result.LastInsertId()
	if err != nil {
		log.Println("Error getting post ID:", err)
		return 0, err
	}

	return int(msgID), nil
}

func GetHistory(userID1, userID2, groupID int, history *[]map[string]interface{}) error {
	var rows *sql.Rows
	var err error

	if groupID == 0 {
		query := `
			SELECT sender_id, content, is_read, created_at
			FROM Messages
			WHERE group_id = 0
			  AND ((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?))
			ORDER BY created_at DESC`
		rows, err = db.Query(query, userID1, userID2, userID2, userID1)
	} else {
		query := `
			SELECT sender_id, content, is_read, created_at
			FROM Messages
			WHERE group_id = ?
			ORDER BY created_at DESC`
		rows, err = db.Query(query, groupID)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var sender int
		var content string
		var isRead bool
		var timestamp string
		if err := rows.Scan(&sender, &content, &isRead, &timestamp); err != nil {
			return err
		}
		username, err := GetUsername(sender)
		if err != nil {
			log.Println("Error fetching username for id: ", sender)
			return err
		}
		message := map[string]interface{}{
			"senderID":       sender,
			"senderUsername": username,
			"content":        content,
			"isRead":         isRead,
			"createdAt":      timestamp,
		}

		*history = append(*history, message)
	}

	if err := rows.Err(); err != nil {
		return err
	}
	return nil

}

func GetMessage(message_id int) ([]string, error) {
	var message []string
	var chatID int
	var senderID int
	var content string
	var createdAt string

	err := db.QueryRow("SELECT chat_id, sender_id, content, created_at FROM Message WHERE id = ?", message_id).Scan(&chatID, &senderID, &content, &createdAt)
	if err != nil {
		return message, err
	}

	username, err := GetUsername(senderID)
	if err != nil {
		log.Println("Error fetching username for id: ", senderID)
		return message, err
	}

	message = []string{
		fmt.Sprint(chatID),
		fmt.Sprint(senderID),
		username,
		content,
		createdAt,
	}

	return message, nil
}
