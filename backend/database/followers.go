package database

// GetFollowersCount counts how many users follow the specified user
func GetFollowersCount(userID int) (int, error) {
	var count int

	// Adjust table and column names to match your schema
	err := db.QueryRow(`
		SELECT COUNT(*) 
		FROM Followers 
		WHERE followed_id = ?`, userID).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetFollowingCount counts how many users the specified user follows
func GetFollowingCount(userID int) (int, error) {
	var count int

	// Adjust table and column names to match your schema
	err := db.QueryRow(`
		SELECT COUNT(*) 
		FROM Followers 
		WHERE follower_id = ?`, userID).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}
