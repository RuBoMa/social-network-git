package database

import (
	"database/sql"
	"social_network/models"
)

// GetFollowersCount counts how many users follow the specified user
func GetFollowersCount(userID int) (int, error) {
	var count int

	// Check how we names the status / active followed
	err := db.QueryRow(`
		SELECT COUNT(*) 
		FROM Followers 
		WHERE followed_id = ? AND status = 'active'`, userID).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetFollowingCount counts how many users the specified user follows
func GetFollowingCount(userID int) (int, error) {
	var count int

	// Check how we names the status / active follower
	err := db.QueryRow(`
		SELECT COUNT(*) 
		FROM Followers 
		WHERE follower_id = ? AND status='active'`, userID).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetFollowers retrieves the list of users that follow the specified user
func GetFollowers(userID int) ([]models.User, error) {
	var followers []models.User

	rows, err := db.Query(`
		SELECT follower_id
		FROM Followers
		WHERE followed_id = ? AND status = 'active'`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var followedID int

		if err := rows.Scan(&followedID); err != nil {
			return nil, err
		}
		followed, err := GetUser(followedID)
		if err != nil {
			return nil, err
		}
		followers = append(followers, followed)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return followers, nil
}

// GetFollowing retrieves the list of users that the specified user follows
func GetFollowing(userID int) ([]models.User, error) {
	var following []models.User

	rows, err := db.Query(`
		SELECT followed_id
		FROM Followers
		WHERE follower_id = ? AND status = 'active'`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int

		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		user, err := GetUser(id)
		if err != nil {
			return nil, err
		}
		following = append(following, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return following, nil
}

// Check if the profile is public or private
func IsProfilePrivate(userID int) (bool, error) {
	var isPublic bool
	query := `SELECT is_public FROM Users WHERE id = ?`
	err := db.QueryRow(query, userID).Scan(&isPublic)
	return !isPublic, err
}

// AddFollower adds a follower to the database or reactivates a deleted relationship
func AddFollower(followerID, followedID int) error {
	var existingStatus string

	// Check if a follower relationship already exists
	queryCheck := `
		SELECT status FROM Followers
		WHERE follower_id = ? AND followed_id = ?
	`
	err := db.QueryRow(queryCheck, followerID, followedID).Scan(&existingStatus)

	if err != nil {
		if err == sql.ErrNoRows {
			// No existing relationship, insert a new one
			insertQuery := `
				INSERT INTO Followers (follower_id, followed_id, status, created_at)
				VALUES (?, ?, 'active', CURRENT_TIMESTAMP)
			`
			_, insertErr := db.Exec(insertQuery, followerID, followedID)
			return insertErr
		}
		// return any DB error
		return err
	}

	// If it exists but was deleted, reactivate it
	if existingStatus == "deleted" {
		updateQuery := `
			UPDATE Followers
			SET status = 'active', updated_at = CURRENT_TIMESTAMP
			WHERE follower_id = ? AND followed_id = ?
		`
		_, updateErr := db.Exec(updateQuery, followerID, followedID)
		return updateErr
	}

	return nil
}

// RemoveFollower sets the follower relationship status to 'deleted' and updates the timestamp
func RemoveFollower(followerID, followedID int) error {
	query := `
		UPDATE Followers
		SET status = 'deleted', updated_at = CURRENT_TIMESTAMP
		WHERE follower_id = ? AND followed_id = ?
	`
	_, err := db.Exec(query, followerID, followedID)
	return err
}
