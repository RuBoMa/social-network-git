package server

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// updates HTTP connection to websocket protocol
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (change for production)
	},
}

var (
	clients       = make(map[*websocket.Conn]int) // Map of WebSocket connections -> userID
	broadcast     = make(chan Message)            // Channel for broadcasting messages
	clientsMutex  sync.Mutex                      // Protects access to activeUsers map
	messagesMutex sync.Mutex
)

type Message struct {
	Type      string    `json:"type"`       // "chat", "message", "update_users"
	Sender    User      `json:"sender"`     // Sender
	Receiver  User      `json:"receiver"`   // Receiver
	GroupID   int       `json:"group_id"`   // Group ID
	CreatedAt string    `json:"created_at"` // Timestamp for the message
	Content   string    `json:"content"`    // Chat message
	IsRead    bool      `json:"is_read"`    // If the message is read
	Users     []User    `json:"users"`      // sorted users with userID and username
	History   []Message `json:"history"`    // Message history
}

type UserInteraction struct {
	UserID          int
	Username        string
	LastInteraction string
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Online   bool   `json:"online"` //If user is currently online
}

type RouteInfo struct {
	Page      string
	PostID    int
	SubAction string
	Err       error
}
