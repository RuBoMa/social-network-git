package app

func SendNotifcation() {

	// Range the users that are going to get the notification
	// Fetch the information based on the related_event_id OR related_request_id
	// Fill in type of the notification (follow_request, group invitation, group_request, event)
	// Fill in needed information
}

// CreateNotification creates a new notification for a user
func CreateNotification(userID, relatedEventID *int) error {
	_, err := db.Exec(`
		INSERT INTO Notifications (user_id, is_read, related_event_id, created_at, updated_at)
		VALUES (?, false, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, userID, relatedEventID)
	return err
}
