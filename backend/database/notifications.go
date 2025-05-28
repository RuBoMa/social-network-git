package database

import (
	"database/sql"
	"fmt"
	"log"
	"social_network/models"
	"time"
)

// AddNotificationIntoDB handles event or request based notification logic
func AddNotificationIntoDB(notifType string, request models.Request, event models.Event) ([]int, error) {
	var query string
	var id int
	var receivers []int
	var notificationIDs []int
	log.Println("Adding notification to DB for type:", notifType)

	switch notifType {
	case "follow_request", "group_invite", "join_request":
		query = `
			INSERT INTO Notifications (user_id, type, is_read, related_event_id, related_request_id, created_at)
			VALUES (?, ?, false, 0, ?, ?)
		`
		id = request.RequestID
		receivers = append(receivers, request.Receiver.UserID)

	case "new_event":
		query = `
			INSERT INTO Notifications (user_id, type, is_read, related_request_id, related_event_id, created_at)
			VALUES (?, ?, false, 0, ?, ?)
		`
		id = event.EventID
		log.Println("Group members for new event:", event.Group.GroupMembers)
		for _, member := range event.Group.GroupMembers {
			if member.UserID != event.Creator.UserID {
				receivers = append(receivers, member.UserID)
			}
		}

	default:
		return notificationIDs, fmt.Errorf("invalid notification type")
	}

	for _, receiver := range receivers {
		result, err := db.Exec(query, receiver, notifType, id, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			log.Println("Error inserting notification into database:", err)
			return notificationIDs, err
		}
		notificationID, err := result.LastInsertId()
		if err != nil {
			log.Println("Error getting last insert ID for notification:", err)
			return notificationIDs, err
		}
		log.Printf("Notification added with ID %d for user %d\n", notificationID, receiver)
		notificationIDs = append(notificationIDs, int(notificationID))

	}
	return notificationIDs, nil
}

func GetNotificationByID(notificationID int) (models.Notification, error) {
	var n models.Notification

	err := db.QueryRow(`
		SELECT id, user_id, type, is_read, related_request_id, related_event_id, created_at
		FROM Notifications
		WHERE id = ?
	`, notificationID).Scan(&n.NotificationID, &n.UserID, &n.Type, &n.IsRead, &n.Request.RequestID, &n.Event.EventID, &n.CreatedAt)
	if err != nil {
		log.Println("Error fetching notification by ID:", err)
		return n, err
	}

	if n.Event.EventID != 0 {
		n.Event, err = GetEventByID(n.Event.EventID)
		if err != nil {
			log.Println("Error fetching event by ID:", err)
			return n, err
		}
		n.Event.Group, err = GetGroupByID(n.Event.Group.GroupID)
		if err != nil {
			log.Println("Error fetching group by ID:", err)
			return n, err
		}
	}
	if n.Request.RequestID != 0 {
		n.Request, err = GetRequestByID(n.Request.RequestID)
		if err != nil {
			log.Println("Error fetching request by ID:", err)
			return n, err
		}
	}

	return n, nil
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

func GetUnreadNotifications(userID int) ([]models.Notification, error) {
	var notifications []models.Notification

	rows, err := db.Query(`
		SELECT id, type, is_read, related_request_id, related_event_id, created_at
		FROM Notifications
		WHERE user_id = ? AND is_read = false
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No unread notifications found for user:", userID)
			return notifications, nil // No unread notifications
		}
		log.Println("Error fetching unread notifications:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var n models.Notification
		if err := rows.Scan(&n.NotificationID, &n.Type, &n.IsRead, &n.Request.RequestID, &n.Event.EventID, &n.CreatedAt); err != nil {
			log.Println("Error scanning notification:", err)
			return nil, err
		}

		if n.Event.EventID != 0 {
			n.Event, err = GetEventByID(n.Event.EventID)
			if err != nil {
				log.Println("Error fetching event by ID:", err)
				return nil, err
			}
			n.Event.Group, err = GetGroupByID(n.Event.Group.GroupID)
			if err != nil {
				log.Println("Error fetching group by ID:", err)
				return nil, err
			}
		}
		if n.Request.RequestID != 0 {
			n.Request, err = GetRequestByID(n.Request.RequestID)
			if err != nil {
				log.Println("Error fetching request by ID:", err)
				return nil, err
			}
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}

// IsValidNotificationID checks if a notification ID exists in the database
// Returns true if the notification ID is valid, false otherwise
func IsValidNotificationID(notificationID int) (bool, error) {
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM Notifications WHERE id = ?)
	`, notificationID).Scan(&exists)
	if err != nil {
		log.Println("Error checking notification ID:", err)
		return false, err
	}
	return exists, nil
}

// NotificationSeen marks a notification as read
func NotificationSeen(notificationID int) error {
	_, err := db.Exec(`UPDATE Notifications SET is_read = true, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, notificationID)
	if err != nil {
		log.Println("Error marking notification as seen:", err)
	}
	return err
}
