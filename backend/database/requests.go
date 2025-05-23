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
	var err error
	if request.Group.GroupID != 0 {
		var id int
		if request.Status == "requested" {
			id = request.Sender.UserID
		} else {
			id = request.Receiver.UserID
		}

		err = db.QueryRow(`
			SELECT id, status FROM Requests
			WHERE ((sent_id = ? AND status = "requested") OR (received_id = ? AND status = "invited"))
			AND group_id = ?
		`,
			id,
			id,
			request.Group.GroupID,
		).Scan(&existingID, &currentStatus)
	} else {
		err = db.QueryRow(`
			SELECT id, status FROM Requests
			WHERE (sent_id = ? AND status = "requested")
		`,
			request.Sender.UserID,
		).Scan(&existingID, &currentStatus)
	}

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
func ActiveRequest(userID, groupID int) (string, int, error) {
	var id int
	var status string
	err := db.QueryRow(`
		SELECT id, status
		FROM requests
		WHERE group_id = ?
		AND (
			(status = 'invited' AND received_id = ?)
			OR
			(status = 'requested' AND sent_id = ?)
		)
		LIMIT 1
	`, groupID, userID, userID).Scan(&id, &status)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No active request found")
			return "", id, nil
		}
		log.Println("Error checking active request:", err)
		return status, id, err
	}
	return status, id, nil
}

func GetRequestByID(requestID int) (models.Request, error) {
	var request models.Request
	err := db.QueryRow(`
		SELECT id, sent_id, received_id, group_id, status, created_at
		FROM Requests
		WHERE id = ?
	`, requestID).Scan(&request.RequestID, &request.Sender.UserID, &request.Receiver.UserID, &request.Group.GroupID, &request.Status, &request.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No active request found")
			return request, nil
		}
		log.Println("Error checking active request:", err)
		return request, err
	}

	request.Group, err = GetGroupByID(request.Group.GroupID)
	if err != nil {
		log.Println("Error getting group info:", err)
		return request, err
	}

	request.Sender, err = GetUser(request.Sender.UserID)
	if err != nil {
		log.Println("Error getting sender info:", err)
		return request, err
	}

	// receiver, err := GetUser(request.Receiver.UserID)
	// if err != nil {
	// 	log.Println("Error getting receiver info:", err)
	// 	return request, err
	// }
	// request.Receiver = receiver

	return request, nil
}

// GetGroupRequests retrieves all requests for a specific group
func GetGroupRequests(groupID int) ([]models.User, error) {
	var users []models.User
	rows, err := db.Query(`
		SELECT sent_id
		FROM Requests
		WHERE group_id = ? AND status = 'requested'
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserID); err != nil {
			return nil, err
		}
		user, err := GetUser(user.UserID)
		if err != nil {
			log.Println("Error getting user info:", err)
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetGroupRequestStatus retrieves a possible request for a user in a group
func GetGroupRequestStatus(groupID, userID int) (models.Request, error) {
	var request models.Request
	err := db.QueryRow(`
		SELECT id, sent_id, status
		FROM Requests
		WHERE group_id = ? AND sent_id = ?
	`, groupID, userID).Scan(&request.RequestID, &request.Sender.UserID, &request.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No active request found")
			return request, nil
		}
		log.Println("Error checking active request:", err)
		return request, err
	}

	return request, nil
}
