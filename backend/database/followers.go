package database

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

// GetFollowing retrieves the list of users that the specified user follows
func GetFollowing(userID int) ([]int, error) {
	var followers []int

	rows, err := db.Query(`
		SELECT follower_id
		FROM Followers
		WHERE followed_id = ? AND status = 'active'`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var followerID int
		if err := rows.Scan(&followerID); err != nil {
			return nil, err
		}
		followers = append(followers, followerID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return followers, nil
}

// Check if the profile is public or private
func IsProfilePrivate(userID int) (bool, error) {
	var isPublic bool
	query := `SELECT is_public FROM Users WHERE id = ?`
	err := db.QueryRow(query, userID).Scan(&isPublic)
	return !isPublic, err
}

// AddFollower adds a follower to the database
func AddFollower(followerID, followedID int) error {
	query := `INSERT INTO Followers (follower_id, followed_id) VALUES (?, ?)`
	_, err := db.Exec(query, followerID, followedID)
	if err != nil {
		return err
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
