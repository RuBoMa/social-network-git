package database

import (
	"database/sql"
	"social_network/models"
	"strings"
	"time"
)

// AddUserIntoDB inserts the user's details into the database
// It takes the user's email, hashed password, first name, last name, date of birth, avatar path, username, about me section, and public status as parameters
func AddUserIntoDB(email, hashedPassword, firstname, lastname, dob, avatar_path, nickname, about_me string, is_public bool) error {

	_, err := db.Exec("INSERT INTO Users (email, password_hash, first_name, last_name, date_of_birth, avatar_path, nickname, about_me, is_public, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		email, hashedPassword, firstname, lastname, dob, avatar_path, nickname, about_me, is_public, time.Now().Format("2006-01-02 15:04:05"))
	return err
}

// isEmailUnique checks if the given email is unique in the database
func IsEmailUnique(email string) (bool, error) {
	email = strings.ToLower(email)

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Users WHERE email = ?", email).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

// getUserCredentials retrieves the user's ID and hashed password from the database
func GetUserCredentials(email string) (int, string, error) {
	var userID int
	var hashedPassword string

	err := db.QueryRow("SELECT id, password_hash FROM Users WHERE email = ?", email).Scan(&userID, &hashedPassword)
	if err != nil {
		return 0, "", err
	}
	return userID, hashedPassword, nil
}

func GetUsers() ([]models.User, error) {

	var users []models.User
	rows, err := db.Query("SELECT id, first_name, last_name, avatar_path, nickname FROM Users WHERE id != 1")
	if err != nil {
		if err == sql.ErrNoRows {
			// No active users, return an empty slice
			return users, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.UserID,
			&user.FirstName,
			&user.LastName,
			&user.AvatarPath,
			&user.Nickname,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func GetActiveUsers() (map[int]string, error) {
	var activeSessions []int
	var activeUsers = make(map[int]string)

	rows, err := db.Query("SELECT user_id FROM Sessions WHERE status = 'active'")
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

	var nickname string
	err := db.QueryRow("SELECT nickname FROM Users WHERE id = ?", userID).Scan(&nickname)
	if err != nil {
		return "", err
	}
	return nickname, nil
}

// GetUser retrieves a user's profile information from the database
func GetUser(userID int) (models.User, error) {
	var user models.User

	err := db.QueryRow(`
        SELECT id, first_name, last_name, avatar_path, nickname, email, date_of_birth, about_me, is_public 
        FROM Users 
        WHERE id = ?`, userID).Scan(
		&user.UserID,
		&user.FirstName,
		&user.LastName,
		&user.AvatarPath,
		&user.Nickname,
		&user.Email,
		&user.DateOfBirth,
		&user.AboutMe,
		&user.IsPublic,
	)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// IsValidUserID checks if the given user ID exists in the database
func IsValidUserID(userID int) bool {
	if userID == 0 {
		return false
	}
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Users WHERE id = ?", userID).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}
