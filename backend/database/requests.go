package database

import (
	"log"
	"social_network/models"
	"time"
)

// AddRequestIntoDB adds a new request to the database
func AddRequestIntoDB(request models.Request) (int, error) {

	var existingID int
	err := db.QueryRow(`
    SELECT id FROM Requests
    WHERE sent_id = ? AND received_id = ? AND group_id = ? AND status = ?
	`, request.Sender.UserID, request.Receiver.UserID, request.Group.GroupID, request.Status).Scan(&existingID)

	if err == nil {
		log.Println("Request already exists with ID:", existingID)
		return existingID, nil
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
