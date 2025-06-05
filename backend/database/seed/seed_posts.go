package seed

import (
	"database/sql"
	"fmt"
	"time"
)

type Post struct {
	ID      int
	UserID  int
	GroupID int
	Title   string
	Content string
	Privacy string
}

func SeedPosts(db *sql.DB) error {

	posts := []Post{
		{ID: 1,
			UserID:  1,
			GroupID: 0,
			Title:   "First Post",
			Content: "This is the content of the first post",
			Privacy: "public"},
		{ID: 2,
			UserID:  2,
			GroupID: 0,
			Title:   "Second Post",
			Content: "This is the content of the second post",
			Privacy: "private"},
		{ID: 3,
			UserID:  3,
			GroupID: 0,
			Title:   "Third Post",
			Content: "This is the content of the third post",
			Privacy: "public"},
	}

	for _, post := range posts {
		_, err := db.Exec("INSERT INTO posts (id, user_id, group_id, title, content, privacy, image_path, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", post.ID, post.UserID, post.GroupID, post.Title, post.Content, post.Privacy, "", time.Now())
		if err != nil {
			return fmt.Errorf("failed to seed posts: %w", err)
		}
	}

	return nil
}
