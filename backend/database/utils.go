package database

import (
	"database/sql"
	"fmt"
)

func GetTimestamp(message_id int, table string) (string, error) {
	var timestamp string

	query := fmt.Sprintf("SELECT created_at FROM %s WHERE id = ?", table)

	err := db.QueryRow(query, message_id).Scan(&timestamp)
	if err != nil {
		return "", err
	}
	return timestamp, nil
}

func GetLastAction(user1, user2 int) (string, error) {
	var timestamp string
	var chatID int

	err := db.QueryRow(
		`SELECT id
		FROM Chat
		WHERE 
			(user1_id = ? AND user2_id = ?) OR
			(user1_id = ? AND user2_id = ?)
		`, user1, user2, user2, user1).Scan(&chatID)

	if err != nil {
		if err == sql.ErrNoRows {
			return timestamp, nil
		} else {
			return timestamp, err
		}
	}

	err = db.QueryRow(
		`SELECT created_at
		FROM Message
		WHERE 
			chat_id = ?
		ORDER BY created_at DESC
		LIMIT 1`, chatID).Scan(&timestamp)

	if err != nil {
		if err == sql.ErrNoRows {
			return timestamp, nil
		} else {
			return timestamp, err
		}
	}
	return timestamp, nil
}
