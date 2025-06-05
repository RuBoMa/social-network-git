package seed

import (
	"database/sql"
	"time"
)

type Message struct {
	ID         int
	SenderID   int
	ReceiverID int
	GroupID    int
	Content    string
	IsRead     bool
}

func SeedMessages(db *sql.DB) error {

	messages := []Message{
		{ID: 1, SenderID: 1, ReceiverID: 0, GroupID: 1, Content: "Hello everyone!", IsRead: true},
		{ID: 2, SenderID: 2, ReceiverID: 0, GroupID: 1, Content: "Hi there!", IsRead: true},
		{ID: 3, SenderID: 1, ReceiverID: 0, GroupID: 1, Content: "What are you doing today?", IsRead: true},
		{ID: 4, SenderID: 2, ReceiverID: 0, GroupID: 1, Content: "Just working on some code.", IsRead: true},
	}

	for _, msg := range messages {
		_, err := db.Exec("INSERT INTO Messages (id, sender_id, received_id, group_id, content, is_read, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
			msg.ID, msg.SenderID, msg.ReceiverID, msg.GroupID, msg.Content, msg.IsRead, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			return err
		}
	}

	return nil
}
