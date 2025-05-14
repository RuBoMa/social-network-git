package database

// UpdateProfileVisibility updates the visibility status of a user's profile
func UpdateProfileVisibility(userID int, isPublic bool) error {
	_, err := db.Exec(`
		UPDATE Users 
		SET is_public = ? 
		WHERE user_id = ?`,
		isPublic, userID)

	return err
}
