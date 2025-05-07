package models

type PostDetails struct {
	PostID      int              `json:"post_id"`
	UserID      int              `json:"user_id"`
	Username    string           `json:"username"`
	PostTitle   string           `json:"post_title"`
	PostContent string           `json:"post_content"`
	Comments    []CommentDetails `json:"comments"`
	PostImage   string           `json:"post_image"`
	CreatedAt   string           `json:"created_at"`
}
type CommentDetails struct {
	CommentID      int    `json:"comment_id"`
	PostID         int    `json:"post_id"`
	UserID         int    `json:"user_id"`
	Username       string `json:"username"`
	CommentContent string `json:"comment_content"`
	CommentImage   string `json:"comment_image"`
	CreatedAt      string `json:"created_at"`
}

// Profile details

// Group details

// Group members

// Event details

// Struct to map the incoming login data
type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Struct to map json response messages NEEDS MORE THAN A MESSAGE ??
type Response struct {
	Message string `json:"message"`
}

// Struct to map the incoming sign-up data
type SignUpData struct {
	Username        string `json:"username"`
	Age             string `json:"age"`
	Gender          string `json:"gender"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}
