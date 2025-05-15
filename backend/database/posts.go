package database

import (
	"log"
	"social_network/models"
	"time"
)

// AddPostToDatabase inserts a new post into the database
// It takes the post title, content, image path, privacy setting, user ID, and group ID as parameters
func AddPostIntoDB(title, content, imagePath, privacy string, userID, groupID int) error {

	_, err := db.Exec("INSERT INTO Posts (user_id, group_id, title, content, image_path, privacy, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		userID, groupID, title, content, imagePath, privacy, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Println("Error inserting post:", err)
		return err
	}

	return nil
}

// AddCommentIntoDB inserts a new comment into the database
// It takes the post ID, user ID, content, and image path as parameters
func AddCommentIntoDB(postID, userID int, content, image_path string) error {
	_, err := db.Exec("INSERT INTO Comments (post_id, user_id, content, image_path, created_at) VALUES (?, ?, ?, ?, ?)",
		postID, userID, content, image_path, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Println("Error creating post:", err)
		return err
	}
	return nil
}

// GetPosts retrieves all posts from the database for a given user.
// It includes public posts, own posts, posts from authors user is following, and custom privacy posts.
func GetPosts(userID int) ([]models.Post, error) {
	var posts []models.Post

	// Query to get all posts ordered by creation date
	query := `
		SELECT Post.id
		FROM Posts AS Post
		JOIN Users ON Post.user_id = Users.id
		LEFT JOIN Followers ON Followers.followed_id = Post.user_id
		LEFT JOIN Post_Privacy ON Post_Privacy.post_id = Post.id
		WHERE 
			Post.privacy = 'public'
			OR Post.user_id = ?
			OR (Post.privacy = 'followers' AND Followers.follower_id = ?)
			OR (Post.privacy = 'custom' AND Post_Privacy.user_id = ? AND Post_Privacy.status = 'active')
		GROUP BY Post.id
		ORDER BY Post.created_at DESC;
	`
	rows, err := db.Query(query, userID, userID, userID)
	if err != nil {
		log.Println("Error fetching posts:", err)
		return nil, err
	}
	defer rows.Close()

	// Loop through the rows and fetch details for each post
	for rows.Next() {
		var postID int
		if err := rows.Scan(&postID); err != nil {
			log.Println("Error scanning post ID:", err)
			return nil, err
		}

		// Get the details for each post
		post, err := GetPostDetails(postID)
		if err != nil {
			log.Println("Error getting post details:", err)
			return nil, err
		}

		posts = append(posts, *post)
	}

	// Return the list of posts
	return posts, nil
}

// GetUserPosts retrieves posts made by a specific user, applying privacy filters based on the viewer
// If isOwnProfile is true, it returns all posts regardless of privacy settings
func GetUserPosts(profileID int, viewerID int, isOwnProfile bool) ([]models.Post, error) {
	var posts []models.Post
	var query string
	var args []interface{}

	if isOwnProfile {
		// If viewing own profile, get all posts regardless of privacy
		query = `
            SELECT Post.id 
            FROM Posts AS Post 
            WHERE Post.user_id = ? 
            ORDER BY Post.created_at DESC
        `
		args = []interface{}{profileID}
	} else {
		// Otherwise, apply privacy filters similar to GetPosts
		query = `
            SELECT Post.id 
            FROM Posts AS Post
            LEFT JOIN Followers ON Followers.followed_id = Post.user_id
            LEFT JOIN Post_Privacy ON Post_Privacy.post_id = Post.id
            WHERE Post.user_id = ? AND (
                Post.privacy = 'public' OR
                (Post.privacy = 'followers' AND Followers.follower_id = ?) OR
                (Post.privacy = 'custom' AND Post_Privacy.user_id = ? AND Post_Privacy.status = 'active')
            )
            GROUP BY Post.id
            ORDER BY Post.created_at DESC
        `
		args = []interface{}{profileID, viewerID, viewerID}
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Println("Error fetching user posts:", err)
		return nil, err
	}
	defer rows.Close()

	// Loop through the rows and fetch details for each post
	for rows.Next() {
		var postID int
		if err := rows.Scan(&postID); err != nil {
			log.Println("Error scanning post ID:", err)
			return nil, err
		}

		// Reuse your existing GetPostDetails function
		post, err := GetPostDetails(postID)
		if err != nil {
			log.Println("Error getting post details:", err)
			return nil, err
		}

		posts = append(posts, *post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// GetPostDetails fetches the details of a specific post from the database
func GetPostDetails(postID int) (*models.Post, error) {

	query := `
		SELECT 
			Post.id AS post_id,
			Post.user_id AS user_id,
			Post.group_id AS group_id,
			Post.title AS post_title,
			Post.content AS post_content,
			Post.image_path AS post_image,
			Post.privacy AS post_privacy,
			Post.created_at AS post_created_at,
			Users.nickname,
			Users.first_name,
			Users.last_name,
			Users.avatar_path
		FROM Posts AS Post
		LEFT JOIN Users ON Post.user_id = Users.id
		WHERE Post.id = ?;
	`

	row := db.QueryRow(query, postID)
	// Scan the data into a PostDetails struct
	post := models.Post{}
	author := models.User{}
	err := row.Scan(
		&post.PostID,
		&author.UserID,
		&post.Group.GroupID,
		&post.PostTitle,
		&post.PostContent,
		&post.PostImage,
		&post.Privacy,
		&post.CreatedAt,
		&author.Nickname,
		&author.FirstName,
		&author.LastName,
		&author.AvatarPath,
	)

	if err != nil {
		log.Println("Error scanning rows")
		return nil, err
	}

	post.Author = author

	post.Comments, err = GetComments(postID)
	if err != nil {
		log.Println("Error getting comments")
		return nil, err
	}

	return &post, nil
}

// GetComments fetches all comments for a specific post from the database
func GetComments(postID int) ([]models.Comment, error) {

	query := `
		SELECT 
			Comment.id AS comment_id,
			Comment.post_id AS post_id,
			Comment.user_id,
			Comment.content AS comment_content,
			Comment.image_path AS comment_image,
			Comment.created_at,
			Users.nickname,
			Users.first_name,
			Users.last_name,
			Users.avatar_path
		FROM Comments AS Comment
		LEFT JOIN Users ON Comment.user_id = Users.id
		WHERE Comment.post_id = ?
		GROUP BY Comment.id, Users.id
		ORDER BY Comment.created_at DESC;
	`
	rows, err := db.Query(query, postID)
	if err != nil {
		log.Println("Error fetching comments from database")
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		comment := models.Comment{}
		author := models.User{}

		err := rows.Scan(
			&comment.CommentID,
			&comment.PostID,
			&author.UserID,
			&comment.CommentContent,
			&comment.CommentImage,
			&comment.CreatedAt,
			&author.Nickname,
			&author.FirstName,
			&author.LastName,
			&author.AvatarPath,
		)
		if err != nil {
			log.Println("Error scanning rows")
			return nil, err
		}

		comment.CommentAuthor = author

		comments = append(comments, comment)
	}

	return comments, nil
}

// ValidatePostID checks if a post with the given ID exists in the database
func ValidatePostID(postID int) bool {
	var post int
	err := db.QueryRow("SELECT id FROM Posts WHERE id = ?", postID).Scan(&post)
	if err != nil {
		log.Println("Error scanning postID:", err)
		return false
	}
	return true
}

// ValidateCommentID checks if a comment with the given ID exists in the database
func ValidateCommentID(commentID int) bool {
	var comment int
	err := db.QueryRow("SELECT id FROM Comments WHERE id = ?", commentID).Scan(&comment)
	if err != nil {
		log.Println("Error scanning commentID:", err)
		return false
	}
	return true
}

// RemoveFromPostPrivacy sets the status to 'inactive' for custom privacy posts
func RemoveFromPostPrivacy(followerID, followedID int) error {
	query := `
		UPDATE Post_Privacy
		SET status = 'inactive', updated_at = CURRENT_TIMESTAMP
		WHERE user_id = ?
		  AND post_id IN (
			  SELECT post_id FROM Post_Privacy
			  WHERE user_id = ? AND status = 'custom'
		  )
	`
	_, err := db.Exec(query, followerID, followedID)
	return err
}
