package database

import (
	"database/sql"
	"log"
	"social_network/models"
	"time"
)

// AddRequestIntoDB adds a new request to the database
func AddRequestIntoDB(request models.Request) (int, error) {

	var existingID int
	var currentStatus string
	err := db.QueryRow(`
    SELECT id, status FROM Requests
    WHERE sent_id = ? AND received_id = ? AND group_id = ?
	`, request.Sender.UserID, request.Receiver.UserID, request.Group.GroupID).Scan(&existingID, &currentStatus)

	if err == nil {
		log.Println("Request already exists with ID:", existingID)
		if currentStatus == request.Status {
			return existingID, nil
		} else {
			_, err := db.Exec("UPDATE Requests SET status = ?, updated_at = ? WHERE id = ?", request.Status, time.Now().Format("2006-01-02 15:04:05"), existingID)
			if err != nil {
				return 0, err
			}
			log.Printf("Request status updated from %s to %s\n", currentStatus, request.Status)
			return existingID, nil
		}
	} else {

		result, err := db.Exec("INSERT INTO Requests (sent_id, received_id, group_id, status, created_at) VALUES (?, ?, ?, ?, ?)",
			request.Sender.UserID, request.Receiver.UserID, request.Group.GroupID, request.Status, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			return 0, err
		}

		invitationID, err := result.LastInsertId()
		if err != nil {
			return 0, err
		}

		return int(invitationID), nil
	}
}

// UpdateRequestStatus updates the status of a request in the database
func UpdateRequestStatus(request models.Request) error {

	_, err := db.Exec("UPDATE Requests SET status = ?, updated_at = ? WHERE id = ?", request.Status, time.Now().Format("2006-01-02 15:04:05"), request.RequestID)
	if err != nil {
		return err
	}
	return nil
}

func IsValidRequestID(requestID int) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Requests WHERE id = ?", requestID).Scan(&count)
	if err != nil {
		log.Println("Error checking request ID:", err)
		return false
	}
	return count > 0
}

// ActiveRequest checks if there is an active request for a user in a group (invitation or own request)
func ActiveRequest(userID, groupID int) (string, error) {
	var status string
	err := db.QueryRow(`
		SELECT status
		FROM Requests
		WHERE (sent_id = ?) AND group_id = ?
	`, userID, groupID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No active request found")
			return "", nil
		}
		log.Println("Error checking active request:", err)
		return status, err
	}
	return status, nil
}

// GetGroupRequests retrieves all requests for a specific group
func GetGroupRequests(groupID int) ([]models.Request, error) {
	var requests []models.Request
	rows, err := db.Query(`
		SELECT id, sent_id
		FROM Requests
		WHERE group_id = ? AND status = 'requested'
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var request models.Request
		if err := rows.Scan(&request.RequestID, &request.Sender.UserID); err != nil {
			return nil, err
		}
		user, err := GetUser(request.Sender.UserID)
		if err != nil {
			log.Println("Error getting user info:", err)
			return nil, err
		}
		request.Sender = user
		requests = append(requests, request)
	}

	return requests, nil
}
