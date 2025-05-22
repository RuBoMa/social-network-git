package models

type Post struct {
	PostID      int       `json:"post_id"`
	Author      User      `json:"author"`
	Group       Group     `json:"group"`
	PostTitle   string    `json:"post_title"`
	PostContent string    `json:"post_content"`
	Comments    []Comment `json:"comments"`
	PostImage   string    `json:"post_image"`
	Privacy     string    `json:"privacy"`
	CustomUsers []User    `json:"custom_users"` // Users tagged in the post
	CreatedAt   string    `json:"created_at"`
}
type Comment struct {
	CommentID      int    `json:"comment_id"`
	PostID         int    `json:"post_id"`
	CommentAuthor  User   `json:"comment_author"`
	CommentContent string `json:"comment_content"`
	CommentImage   string `json:"comment_image"`
	CreatedAt      string `json:"created_at"`
}
