package database

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"
)

// insertUserIntoDB inserts the user's details into the database
func InsertUserIntoDB(username, age, gender, firstname, lastname, email, hashedPassword string) error {
	ageInt, err := strconv.Atoi(age)
	if err != nil {
		ageInt = 0
	}
	_, err = db.Exec("INSERT INTO Users (username, age, gender, firstname, lastname, email, password, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		username, ageInt, gender, firstname, lastname, email, hashedPassword, time.Now().Format("2006-01-02 15:04:05"))
	return err
}

// AddPostToDatabase inserts a new post into the database
func AddPostToDatabase(title, content string, userID int) error {

	var err error
	_, err = db.Exec("INSERT INTO Post (title, content, user_id, created_at) VALUES (?, ?, ?, ?)",
		title, content, userID, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Println("Error inserting post:", err)
		return err
	}

	return nil
}

func AddMessageToDB(userID int, content string, chatID int) (int, error) {

	var result sql.Result
	var err error
	result, err = db.Exec("INSERT INTO Message (chat_id, sender_id, content, created_at) VALUES (?, ?, ?, ?)",
		chatID, userID, content, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Println("Error inserting post:", err)
		return 0, err
	}

	// Get the post id for the post inserted
	msgID, err := result.LastInsertId()
	if err != nil {
		log.Println("Error getting post ID:", err)
		return 0, err
	}

	return int(msgID), nil
}

// AddVotes adds or updates a vote type for a post or comment
func AddVotes(userID, postID, commentID, vote int) error {
	var row *sql.Row
	query := `SELECT Type FROM Like WHERE user_id = ? AND `
	deleteQuery := `UPDATE Like SET type = 0, updated_at = ? WHERE user_id = ? AND `
	updateQuery := `UPDATE Like SET type = ?, updated_at = ? WHERE user_id = ? AND `
	var addon string
	var ID int

	if postID == 0 && commentID == 0 {
		return fmt.Errorf("both postID and commentID cannot be zero")
	}

	if postID == 0 {
		ID = commentID
		addon = `comment_id = ?`
	} else if commentID == 0 {
		ID = postID
		addon = `post_id = ?`
	}
	query += addon
	deleteQuery += addon
	updateQuery += addon
	// Check if the user has already liked the post or comment
	row = db.QueryRow(query, userID, ID)
	var likeType int
	err := row.Scan(&likeType)
	if err != nil {
		if err == sql.ErrNoRows {
			likeType = -1 // To imply that no record exists
		} else {
			log.Println("Error scanning current value:", err)
			return err
		}
	}

	if likeType == vote {
		// If existing like type is the same the the current, remove the like by changing the type to 0
		_, err = db.Exec(deleteQuery, time.Now().Format("2006-01-02 15:04:05"), userID, ID)
		if err != nil {
			log.Println("Error updating the record to 0:", err)
			return err
		}
	} else if likeType == -1 {
		// If no record exists, insert a new one
		insertQuery := `INSERT INTO Like (type, user_id, post_id, comment_id, created_at) VALUES (?, ?, ?, ?, ?)`
		_, err = db.Exec(insertQuery, vote, userID, postID, commentID, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			log.Println("Error inserting record:", err)
			return err
		}
	} else {
		_, err = db.Exec(updateQuery, vote, time.Now().Format("2006-01-02 15:04:05"), userID, ID)
		if err != nil {
			log.Println("Error updating the record to new vote:", err)
			return err
		}
	}
	return nil
}

func AddComment(postID int, content string, userID int) error {
	_, err := db.Exec("INSERT INTO Comment (post_id, content, user_id, created_at) VALUES (?, ?, ?, ?)",
		postID, content, userID, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Println("Error creating post:", err)
		return err
	}
	return nil
}
