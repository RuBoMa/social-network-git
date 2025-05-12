package models

// * -- Chat structs --*
type ChatMessage struct {
	Type      string        `json:"type"`       // "chat", "message", "update_users"
	Sender    User          `json:"sender"`     // Sender
	Receiver  User          `json:"receiver"`   // Receiver
	GroupID   int           `json:"group_id"`   // Group ID
	CreatedAt string        `json:"created_at"` // Timestamp for the message
	Content   string        `json:"content"`    // Chat message
	IsRead    bool          `json:"is_read"`    // If the message is read
	Users     []User        `json:"users"`      // sorted users with userID and username
	History   []ChatMessage `json:"history"`    // Message history
}

type UserInteraction struct {
	User            User
	LastInteraction string
}
