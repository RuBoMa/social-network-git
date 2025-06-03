package database

import (
	"database/sql"
	"log"
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

// IsEmailAndNicknameUnique checks if the given email and nickname are unique in the database
func IsEmailAndNicknameUnique(email, nickname string) (bool, bool, error) {
	email = strings.ToLower(email)
	uniqueNickname := false
	uniqueEmail := false

	var emailCount, nicknameCount int
	err := db.QueryRow("SELECT COUNT(*) FROM Users WHERE email = ?", email).Scan(&emailCount)
	if err != nil {
		return false, false, err
	}
	if nickname != "" {
		err = db.QueryRow("SELECT COUNT(*) FROM Users WHERE nickname = ?", nickname).Scan(&nicknameCount)
		if err != nil {
			return false, false, err
		}
	}

	if emailCount == 0 {
		uniqueEmail = true
	}
	if nicknameCount == 0 || nickname == "" {
		uniqueNickname = true
	}

	return uniqueEmail, uniqueNickname, nil
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
	rows, err := db.Query("SELECT id, first_name, last_name, avatar_path, nickname FROM Users")
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

// SearchUsers retrieves max 10 users from the database based on a search term
// It returns a slice of User structs with basic information about the user that match the search criteria
func SearchUsers(searchTerm string, userID int) ([]models.User, error) {
	var users []models.User

	rows, err := db.Query(`
		SELECT id, first_name, last_name, avatar_path, nickname
		FROM Users
		WHERE (nickname LIKE ? OR first_name LIKE ? OR last_name LIKE ?) AND id != ?
		ORDER BY
		CASE
			WHEN nickname = ? OR first_name = ? OR last_name = ? THEN 0
			ELSE 1
		END,
		CASE
			WHEN nickname IS NOT NULL AND nickname != '' THEN nickname
			ELSE last_name
		END ASC
		LIMIT 10
	`,
		"%"+searchTerm+"%", "%"+searchTerm+"%", "%"+searchTerm+"%", userID, // for partial matches
		searchTerm, searchTerm, searchTerm, // for exact matches
	)

	if err != nil {
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

func GetNonGroupMembers(groupID int) ([]models.User, error) {
	var users []models.User

	rows, err := db.Query(`
		SELECT id, first_name, last_name, avatar_path, nickname
		FROM Users
		WHERE id NOT IN (
			SELECT user_id FROM Group_Members WHERE group_id = ?
		)
		AND id NOT IN (
			SELECT joining_user_id FROM requests WHERE group_id = ? AND status IN ('requested')
		)
		AND id NOT IN (
			SELECT joining_user_id FROM requests WHERE group_id = ? AND status IN ('invited')
		)
			ORDER BY
		CASE
			WHEN nickname IS NOT NULL AND nickname != '' THEN nickname
			ELSE last_name
		END ASC
	`, groupID, groupID, groupID)

	if err != nil {
		log.Println("Error executing query:", err)
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
		log.Println("Error scanning rows:", err)
		return nil, err
	}

	return users, nil
}

// UpdatePrivacySettings updates the user's privacy settings in the database
func UpdatePrivacySettings(userID int, isPublic bool) error {
	_, err := db.Exec("UPDATE Users SET is_public = ?, updated_at = ? WHERE id = ?", isPublic, time.Now().Format("2006-01-02 15:04:05"), userID)
	if err != nil {
		log.Println("Error updating privacy settings:", err)
		return err
	}
	return nil
}

// getinteractedUsers retrieves a list of users with whom the specified user has exchanged private messages
func GetInteractedUsers(userID int) ([]models.User, error) {
	// Replace this query with your real logic to find users with whom this user has exchanged messages
	query := `
		SELECT DISTINCT u.id, u.nickname, u.first_name, u.last_name, u.avatar_path
		FROM users u
		INNER JOIN messages m ON (u.id = m.sender_id OR u.id = m.received_id)
		WHERE (m.sender_id = ? OR m.received_id = ?) AND m.group_id = 0 AND u.id != ?
	`

	rows, err := db.Query(query, userID, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserID, &user.Nickname, &user.FirstName, &user.LastName, &user.AvatarPath); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// GetInteractedGroups retrieves a list of groups that has active chat
func GetInteractedGroups(userID int) ([]models.Group, error) {
	var groups []models.Group

	query := `
        SELECT gm.group_id
        FROM Group_Members gm
        JOIN Messages m ON gm.group_id = m.group_id
        WHERE gm.user_id = ?
          AND m.group_id != 0
        GROUP BY gm.group_id
    `
	rows, err := db.Query(query, userID)
	if err != nil {
		return groups, err
	}
	defer rows.Close()

	for rows.Next() {
		var group models.Group
		if err := rows.Scan(&group.GroupID); err != nil {
			return nil, err
		}
		group, err = GetGroupByID(group.GroupID)
		if err != nil {
			log.Println("Error fetching group details for group ID:", group.GroupID, err)
			continue // Skip this group if there's an error
		}
		groups = append(groups, group)
	}
	return groups, nil
}
