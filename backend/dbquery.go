package backend

import (
	"database/sql"
	"fmt"
	"log"
	"social_network/database"
	"strings"
	"time"
)

// GetCategories retrieves all categories from the database
func GetCategories() ([]CategoryDetails, error) {
	rows, err := db.Query("SELECT id, name FROM Category")
	if err != nil {
		log.Println("Error retrieving categories:", err)
		return nil, err
	}
	defer rows.Close()

	var categories []CategoryDetails
	for rows.Next() {
		var category CategoryDetails
		if err := rows.Scan(&category.CategoryID, &category.CategoryName); err != nil {
			log.Println("Error scanning category:", err)
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

// GetPostDetails fetches the details of a specific post from the database
func GetPostDetails(postID, userID int) (*PostDetails, error) {

	row := db.QueryRow(database.PostContent(), postID)
	var err error
	// Scan the data into a PostDetails struct
	post := PostDetails{}
	var categories string
	err = row.Scan(
		&post.PostID,
		&post.UserID,
		&post.Username,
		&post.PostTitle,
		&post.PostContent,
		&post.CreatedAt,
		&post.Likes,
		&post.Dislikes,
		&categories,
	)

	if err != nil {
		log.Println("Error scanning rows")
		return nil, err
	}

	if categories != "" {
		post.Categories = strings.Split(categories, ",")
	}

	post.Comments, err = GetComments(postID, userID)
	if err != nil {
		log.Println("Error getting comments")
		return nil, err
	}

	post.LikedNow, post.DislikedNow, err = GetLikes(userID, postID, 0)
	if err != nil {
		log.Println("Error getting votes")
		return nil, err
	}

	return &post, nil
}

// GetComments fetches all comments for a specific post from the database
func GetComments(postID, userID int) ([]CommentDetails, error) {

	rows, err := db.Query(database.CommentContent(), postID)
	if err != nil {
		log.Println("Error fetching comments from database")
		return nil, err
	}
	defer rows.Close()

	var comments []CommentDetails
	for rows.Next() {
		var comment CommentDetails
		err := rows.Scan(
			&comment.CommentID,
			&comment.PostID,
			&comment.Content,
			&comment.UserID,
			&comment.CreatedAt,
			&comment.Username,
			&comment.Likes,
			&comment.Dislikes,
		)
		if err != nil {
			log.Println("Error scanning rows")
			return nil, err
		}

		comment.LikedNow, comment.DislikedNow, err = GetLikes(userID, 0, comment.CommentID)
		if err != nil {
			log.Println("Error getting votes")
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

// GetLikes fetches the votes for a specific post or comment from the database
func GetLikes(userID, postID, commentID int) (bool, bool, error) {
	if userID != 0 {
		var rows *sql.Rows
		var err error
		if postID == 0 {
			rows, err = db.Query(database.Likes(), userID, nil, commentID)
		} else {
			rows, err = db.Query(database.Likes(), userID, postID, nil)
		}
		if err != nil {
			log.Println("Error fetching votes from database")
			return false, false, err
		}
		defer rows.Close()

		var voteType int
		for rows.Next() {
			err := rows.Scan(&voteType)
			if err != nil {
				log.Println("Error scanning rows")
				return false, false, err
			}
		}
		if voteType == 1 {
			return true, false, nil
		} else if voteType == 2 {
			return false, true, nil
		}
	}

	return false, false, nil
}

// getUserCredentials retrieves the user's ID and hashed password from the database
func getUserCredentials(username string) (int, string, error) {
	var userID int
	var hashedPassword string

	err := db.QueryRow("SELECT id, password FROM User WHERE username = ?", username).Scan(&userID, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			err = db.QueryRow("SELECT id, password FROM User WHERE email = ?", username).Scan(&userID, &hashedPassword)
			if err != nil {
				return 0, "", err
			}
		} else {
			return 0, "", err
		}
	}
	return userID, hashedPassword, nil
}

func GetUsers() (map[int]string, error) {
	var users = make(map[int]string)

	rows, err := db.Query("SELECT id, username FROM User WHERE id != 1")
	if err != nil {
		if err == sql.ErrNoRows {
			// No active users, return an empty slice
			return users, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userID int
		var username string
		if err := rows.Scan(&userID, &username); err != nil {
			return nil, err
		}
		users[userID] = username
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func GetActiveUsers() (map[int]string, error) {
	var activeSessions []int
	var activeUsers = make(map[int]string)

	rows, err := db.Query("SELECT user_id FROM Session WHERE status = 'active'")
	if err != nil {
		if err == sql.ErrNoRows {
			// No active users, return an empty slice
			return activeUsers, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		activeSessions = append(activeSessions, userID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for _, user := range activeSessions {
		username, err := GetUsername(user)
		if err != nil {
			return nil, err
		}
		if username != "" {
			activeUsers[user] = username
		}
	}
	return activeUsers, nil
}

func GetUsername(userID int) (string, error) {

	var username string
	err := db.QueryRow("SELECT username FROM User WHERE id = ?", userID).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}

func GetChatID(user1, user2 int) (int, error) {
	var chatID int

	// Query the database for the chat ID, considering both (user1, user2) and (user2, user1)
	query := `
        SELECT id
        FROM Chat
        WHERE 
            (user1_id = ? AND user2_id = ?) OR
            (user1_id = ? AND user2_id = ?)
        LIMIT 1
    `

	// Try to get the chatID if already exists
	err := db.QueryRow(query, user1, user2, user2, user1).Scan(&chatID)
	if err != nil {
		if err == sql.ErrNoRows {
			insertQuery := `
				INSERT INTO Chat (user1_id, user2_id, created_at)
				VALUES (?, ?, ?)
			`
			// Insert the new chat into the database
			res, err := db.Exec(insertQuery, user1, user2, time.Now().Format("2006-01-02 15:04:05"))
			if err != nil {
				return 0, err
			}

			// Get the last inserted ID
			latestID, err := res.LastInsertId()
			chatID = int(latestID)
			if err != nil {
				return 0, err
			}
		} else {
			return 0, err
		}
	}

	return int(chatID), nil
}

func GetParticipants(chatID int) ([]int, error) {
	var user1ID, user2ID int

	err := db.QueryRow("SELECT user1_id, user2_id FROM Chat WHERE id = ?", chatID).Scan(&user1ID, &user2ID)
	if err != nil {
		return nil, err
	}

	return []int{user1ID, user2ID}, nil
}

func GetHistory(chatID int, history *[]map[string]interface{}) error {

	rows, err := db.Query("SELECT sender_id, content, created_at FROM Message WHERE chat_id = ? AND status = 'active' ORDER BY created_at DESC", chatID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var sender int
		var content string
		var timestamp string
		if err := rows.Scan(&sender, &content, &timestamp); err != nil {
			return err
		}
		username, err := GetUsername(sender)
		if err != nil {
			log.Println("Error fetching username for id: ", sender)
			return err
		}
		message := map[string]interface{}{
			"senderID":       sender,
			"senderUsername": username,
			"createdAt":      timestamp,
			"content":        content,
		}

		*history = append(*history, message)
	}

	if err := rows.Err(); err != nil {
		return err
	}
	return nil

}

func GetTimestamp(message_id int, table string) (string, error) {
	var timestamp string

	query := fmt.Sprintf("SELECT created_at FROM %s WHERE id = ?", table)

	err := db.QueryRow(query, message_id).Scan(&timestamp)
	if err != nil {
		return "", err
	}
	return timestamp, nil
}

func GetMessage(message_id int) ([]string, error) {
	var message []string
	var chatID int
	var senderID int
	var content string
	var createdAt string

	err := db.QueryRow("SELECT chat_id, sender_id, content, created_at FROM Message WHERE id = ?", message_id).Scan(&chatID, &senderID, &content, &createdAt)
	if err != nil {
		return message, err
	}

	username, err := GetUsername(senderID)
	if err != nil {
		log.Println("Error fetching username for id: ", senderID)
		return message, err
	}

	message = []string{
		fmt.Sprint(chatID),
		fmt.Sprint(senderID),
		username,
		content,
		createdAt,
	}

	return message, nil
}

func GetLastAction(user1, user2 int) (string, error) {
	var timestamp string
	var chatID int

	err := db.QueryRow(
		`SELECT id
		FROM Chat
		WHERE 
			(user1_id = ? AND user2_id = ?) OR
			(user1_id = ? AND user2_id = ?)
		`, user1, user2, user2, user1).Scan(&chatID)

	if err != nil {
		if err == sql.ErrNoRows {
			return timestamp, nil
		} else {
			return timestamp, err
		}
	}

	err = db.QueryRow(
		`SELECT created_at
		FROM Message
		WHERE 
			chat_id = ?
		ORDER BY created_at DESC
		LIMIT 1`, chatID).Scan(&timestamp)

	if err != nil {
		if err == sql.ErrNoRows {
			return timestamp, nil
		} else {
			return timestamp, err
		}
	}
	return timestamp, nil
}
