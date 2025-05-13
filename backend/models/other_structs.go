package models

// Any requests made (follow requests, group join requests, etc.)
type Request struct {
	RequestID int    `json:"request_id"`
	Sender    User   `json:"sender"`
	Receiver  User   `json:"receiver"`
	Group     Group  `json:"group"`
	Status    string `json:"status"` // "pending", "accepted", "rejected"
	CreatedAt string `json:"created_at"`
}

// Notifications
type Notification struct {
	NotificationID int     `json:"notification_id"`
	UserID         int     `json:"user_id"`
	Type           string  `json:"type"`
	Event          Event   `json:"event"`
	FollowRequest  Request `json:"follow_request"`
	Seen           bool    `json:"seen"`
	CreatedAt      string  `json:"created_at"`
}

// Struct to map the incoming login data
type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Struct to map json response messages NEEDS MORE THAN A MESSAGE ??
type Response struct {
	Message string `json:"message"`
}

// Struct to map the incoming URL data
type RouteInfo struct {
	Page      string
	PostID    int
	SubAction string
	Err       error
}
