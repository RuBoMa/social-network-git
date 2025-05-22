package database

import (
	"fmt"
	"log"
	"social_network/models"
	"time"
)

// CreateNotification inserts a notification into the database
func CreateNotification(notifType string, userID int, requestID, eventID *int) error {
	_, err := db.Exec(`
		INSERT INTO Notifications (user_id, type, is_read, related_request_id, related_event_id, created_at, updated_at)
		VALUES (?, ?, false, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, userID, notifType, requestID, eventID)
	return err
}

// AddNotificationIntoDB handles event or request based notification logic
func AddNotificationIntoDB(notifType string, request models.Request, event models.Event) error {
	var query string
	var id int
	var receivers []int

	switch notifType {
	case "follow_request", "group_invite", "join_request":
		query = `
			INSERT INTO Notifications (user_id, type, is_read, related_request_id, created_at)
			VALUES (?, ?, false, ?, ?)
		`
		id = request.RequestID
		receivers = append(receivers, request.Receiver.UserID)

	case "event_created":
		query = `
			INSERT INTO Notifications (user_id, type, is_read, related_event_id, created_at)
			VALUES (?, ?, false, ?, ?)
		`
		id = event.EventID
		for _, member := range event.Group.GroupMembers {
			if member.UserID != event.Creator.UserID {
				receivers = append(receivers, member.UserID)
			}
		}

	default:
		return fmt.Errorf("invalid notification type")
	}

	for _, receiver := range receivers {
		_, err := db.Exec(query, receiver, notifType, id, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			log.Println("Error inserting notification into database:", err)
			return err
		}
	}
	return nil
}

// GetNotifications fetches all notifications for a user
func GetNotifications(userID int) ([]models.Notification, error) {
	var notifications []models.Notification

	rows, err := db.Query(`
		SELECT id, user_id, type, is_read, related_request_id, related_event_id, created_at
		FROM Notifications
		WHERE user_id = ?
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		log.Println("Error fetching notifications:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var n models.Notification
		if err := rows.Scan(&n.NotificationID, &n.UserID, &n.Type, &n.IsRead, &n.Request.RequestID, &n.Event.EventID, &n.CreatedAt); err != nil {
			log.Println("Error scanning notification:", err)
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, rows.Err()
}

// NotificationSeen marks a notification as read
func NotificationSeen(notificationID int) error {
	_, err := db.Exec(`UPDATE Notifications SET is_read = true, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, notificationID)
	if err != nil {
		log.Println("Error marking notification as seen:", err)
	}
	return err
}
