package database

import (
	"fmt"
	"log"
	"social_network/models"
	"time"
)

func AddNotificationIntoDB(request models.Request, event models.Event) error {
	var query string
	var ID int
	var receivers []int
	if request.RequestID != 0 {
		query = `
		INSERT INTO Notifications (user_id, is_read, related_request_id, created_at)
		VALUES (?, ?, ?, ?)
		`
		ID = request.RequestID
		receivers = append(receivers, request.Receiver.UserID)
	} else if event.EventID != 0 {
		query = `
		INSERT INTO Notifications (user_id, is_read, related_event_id, created_at)
		VALUES (?, ?, ?, ?)
		`
		ID = event.EventID
		for _, member := range event.Group.GroupMembers {
			if member.UserID == event.CreatorID {
				continue
			}
			receivers = append(receivers, member.UserID)
		}

	} else {
		return fmt.Errorf("no valid request or event ID provided")
	}

	for _, receiver := range receivers {
		_, err := db.Exec(query, receiver, false, ID, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			log.Println("Error inserting notification into database:", err)
			return err
		}
	}
	return nil

}

// GetNotifications retrieves notifications for a specific user from the database
// It returns all notifications no matter if they are read or not
func GetNotifications(userID int) ([]models.Notification, error) {
	var notifications []models.Notification

	rows, err := db.Query(`
		SELECT n.id, n.user_id, n.is_read, n.related_request_id, n.related_event_id, n.created_at
		FROM Notifications n
		WHERE n.user_id = ?
		ORDER BY n.created_at DESC`, userID)
	if err != nil {
		log.Println("Error fetching notifications:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var notification models.Notification
		if err := rows.Scan(&notification.NotificationID, &notification.UserID, &notification.IsRead, &notification.Request.RequestID, &notification.Event.EventID, &notification.CreatedAt); err != nil {
			log.Println("Error scanning notification:", err)
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over notifications:", err)
		return nil, err
	}

	return notifications, nil
}

// NotificationSeen updates the is_read column into "true" in the database
func NotificationSeen(notificationID int) error {
	_, err := db.Exec("UPDATE Notifications SET is_read = ? WHERE id = ?", true, notificationID)
	if err != nil {
		log.Println("Error updating notification status:", err)
		return err
	}
	return nil
}
