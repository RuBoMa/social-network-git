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

	err := db.QueryRow(`
		SELECT created_at
		FROM Messages
		WHERE 
			group_id = 0 AND 
			((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?))
		ORDER BY created_at DESC
		LIMIT 1
	`, user1, user2, user2, user1).Scan(&timestamp)

	if err != nil {
		if err == sql.ErrNoRows {
			return timestamp, nil
		} else {
			return timestamp, err
		}
	}
	return timestamp, nil
}
