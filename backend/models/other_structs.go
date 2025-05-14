package models

// Any requests made (follow requests, group join requests, etc.)
// If group is empty, it means that the request is a follow request
// If receiver is empty, it means that the request is a group join request
// If sender, receiver and group exist, it means that the request is a group invitation
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
	Request        Request `json:"request"`
	IsRead         bool    `json:"is_read"`
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
