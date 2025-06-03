package database

import (
	"database/sql"
	"log"
	"social_network/models"
	"time"
)

// AddRequestIntoDB adds a new request to the database
// If the request already exists, it updates the status if it's different.
// It returns the request ID and any error encountered.
func AddRequestIntoDB(request models.Request) (int, error) {

	var existingID int
	var currentStatus string
	var err error

	// Check if the group request already exists
	if request.Group.GroupID != 0 {
		err = db.QueryRow(`
			SELECT id, status FROM Requests
			WHERE joining_user_id = ? AND (status = "requested" OR status = "invited")
			AND group_id = ?
		`,
			request.JoiningUser.UserID,
			request.Group.GroupID,
		).Scan(&existingID, &currentStatus)

	} else { // Check if the follow request already exists
		err = db.QueryRow(`
			SELECT id, status FROM Requests
			WHERE (sent_id = ? AND status = "requested")
		`,
			request.Sender.UserID,
		).Scan(&existingID, &currentStatus)
	}

	// If the request already exists, update it if the status is different
	if err == nil {
		if currentStatus == request.Status {
			return existingID, nil
		} else {
			_, err := db.Exec("UPDATE Requests SET status = ?, updated_at = ? WHERE id = ?", request.Status, time.Now().Format("2006-01-02 15:04:05"), existingID)
			if err != nil {
				return 0, err
			}
			return existingID, nil
		}
	} else if err == sql.ErrNoRows { // If no existing request found, insert a new one

		result, err := db.Exec("INSERT INTO Requests (sent_id, received_id, joining_user_id, group_id, status, created_at) VALUES (?, ?, ?, ?, ?, ?)",
			request.Sender.UserID, request.Receiver.UserID, request.JoiningUser.UserID, request.Group.GroupID, request.Status, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			return 0, err
		}

		requestID, err := result.LastInsertId()
		if err != nil {
			return 0, err
		}

		return int(requestID), nil
	}

	return 0, err
}

// UpdateRequestStatus updates the status of a request in the database
func UpdateRequestStatus(request models.Request) error {

	_, err := db.Exec("UPDATE Requests SET status = ?, updated_at = ? WHERE id = ?", request.Status, time.Now().Format("2006-01-02 15:04:05"), request.RequestID)
	if err != nil {
		return err
	}
	return nil
}

// IsValidRequestID checks if a request ID exists in the database
func IsValidRequestID(requestID int) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Requests WHERE id = ?", requestID).Scan(&count)
	if err != nil {
		log.Println("Error checking request ID:", err)
		return false
	}
	return count == 1
}

// ActiveRequest checks if there is an active request for a user in a group (invitation or own request)
func ActiveGroupRequest(userID, groupID int) (string, int, error) {
	var id int
	var status string
	err := db.QueryRow(`
		SELECT id, status
		FROM requests
		WHERE group_id = ?
		AND joining_user_id = ?
		AND (status = 'invited' OR status = 'requested')
		LIMIT 1
	`, groupID, userID).Scan(&id, &status)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No active request found")
			return status, id, nil
		}
		log.Println("Error checking active request:", err)
		return status, id, err
	}
	return status, id, nil
}

func GetRequestByID(requestID int) (models.Request, error) {
	var request models.Request
	err := db.QueryRow(`
		SELECT id, sent_id, received_id, joining_user_id, group_id, status, created_at
		FROM Requests
		WHERE id = ?
	`, requestID).Scan(&request.RequestID, &request.Sender.UserID, &request.Receiver.UserID, &request.JoiningUser.UserID, &request.Group.GroupID, &request.Status, &request.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No active request found")
			return request, nil
		}
		log.Println("Error checking active request:", err)
		return request, err
	}
	if request.Group.GroupID > 0 {
		request.Group, err = GetGroupByID(request.Group.GroupID)
		if err != nil {
			log.Println("Error getting group info:", err)
			return request, err
		}
	}

	if request.Sender.UserID > 0 {
		request.Sender, err = GetUser(request.Sender.UserID)
		if err != nil {
			log.Println("Error getting sender info:", err)
			return request, err
		}
	}

	if request.JoiningUser.UserID > 0 {
		request.JoiningUser, err = GetUser(request.JoiningUser.UserID)
		if err != nil {
			log.Println("Error getting joining user info:", err)
			return request, err
		}
	}

	return request, nil
}

// GetGroupRequests retrieves all requests for a specific group
func GetGroupRequests(groupID int) ([]models.Request, error) {
	var requests []models.Request
	rows, err := db.Query(`
		SELECT id, joining_user_id
		FROM Requests
		WHERE group_id = ? AND status = 'requested'
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var request models.Request
		if err := rows.Scan(&request.RequestID, &request.JoiningUser.UserID); err != nil {
			return nil, err
		}
		if request.JoiningUser.UserID != 0 {
			request.JoiningUser, err = GetUser(request.JoiningUser.UserID)
			if err != nil {
				log.Println("Error getting user info:", err)
				return nil, err
			}
		}
		requests = append(requests, request)
	}

	return requests, nil
}

// GetGroupRequestStatus retrieves a possible request for a user in a group
func GetGroupRequestStatus(groupID, userID int) (models.Request, error) {
	var request models.Request
	err := db.QueryRow(`
		SELECT id, sent_id, status
		FROM Requests
		WHERE group_id = ? AND joining_user_id = ?
	`, groupID, userID).Scan(&request.RequestID, &request.JoiningUser.UserID, &request.Status)
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

// HasPendingFollowRequest checks if there is already a pending follow request between two users
func HasPendingFollowRequest(senderID, receiverID int) (bool, error) {
	var id int
	err := db.QueryRow(`
		SELECT id FROM Requests
		WHERE sent_id = ? AND received_id = ?
		AND status = 'follow'
		LIMIT 1
	`, senderID, receiverID).Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		log.Println("Error checking for pending follow request:", err)
		return false, err
	}

	return true, nil

}

// GetOwnFollowRequests retrieves all follow requests sent to the user
func GetOwnFollowRequests(userID int) ([]models.Request, error) {
	var requests []models.Request
	rows, err := db.Query(`
		SELECT id, sent_id, created_at
		FROM requests
		WHERE received_id = ? AND status = 'follow'
	`, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No follow requests found")
			return requests, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var request models.Request
		if err := rows.Scan(&request.RequestID, &request.Sender.UserID, &request.CreatedAt); err != nil {
			return nil, err
		}
		request.Sender, err = GetUser(request.Sender.UserID)
		if err != nil {
			log.Println("Error getting user info:", err)
			return nil, err
		}
		requests = append(requests, request)
	}

	return requests, nil

}
