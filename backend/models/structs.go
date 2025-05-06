package models

type PostDetails struct {
	PostID      int              `json:"post_id"`
	UserID      int              `json:"user_id"`
	Username    string           `json:"username"`
	PostTitle   string           `json:"post_title"`
	PostContent string           `json:"post_content"`
	Comments    []CommentDetails `json:"comments"`
	Categories  []string         `json:"categories"`
	CreatedAt   string           `json:"created_at"`
	Likes       int              `json:"likes"`
	Dislikes    int              `json:"dislikes"`
	LikedNow    bool             `json:"liked_now"`
	DislikedNow bool             `json:"disliked_now"`
}
type CommentDetails struct {
	CommentID   int    `json:"comment_id"`
	PostID      int    `json:"post_id"`
	Content     string `json:"comment_content"`
	UserID      int    `json:"user_id"`
	Username    string `json:"username"`
	CreatedAt   string `json:"created_at"`
	Likes       int    `json:"likes"`
	Dislikes    int    `json:"dislikes"`
	LikedNow    bool   `json:"liked_now"`
	DislikedNow bool   `json:"disliked_now"`
}

type VoteDetails struct {
	Vote      string `json:"vote"`
	CommentID int    `json:"comment_id"`
	PostID    int    `json:"post_id"`
}

type CategoryDetails struct {
	CategoryID   int
	CategoryName string
}

// Struct to map the incoming login data
type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Struct to map json response messages
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
