package app

import (
	"fmt"
	"log"
)

func SendNotifcation() {

	// Range the users that are going to get the notification
	// Fetch the information based on the related_event_id OR related_request_id
	// Fill in type of the notification (follow_request, group invitation, group_request, event)
	// Fill in needed information
}
// CreateNotification creates a new notification for a user
func CreateNotification(userID, message string, relatedEventID *int) error {
	_, err := db.Exec(`
		INSERT INTO Notifications (user_id, is_read, related_event_id, created_at, updated_at)
		VALUES (?, false, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, userID, relatedEventID)
	return err
}

// SendNewFollowerNotification sends a notification when a user starts following another user
func SendNewFollowerNotification(followerID, followedID string) {
	message := fmt.Sprintf("User %s started following you", followerID)
	err := CreateNotification(followedID, message, nil)
	if err != nil {
		log.Println("Failed to create follower notification:", err)
	}
}

// SendFollowRequestNotification sends a notification when a user requests to follow another user
func SendFollowRequestNotification(followerID, followedID string) {
	message := fmt.Sprintf("User %s requested to follow you", followerID)
	err := CreateNotification(followedID, message, nil)
	if err != nil {
		log.Println("Failed to create follower notification:", err)
	}
}
