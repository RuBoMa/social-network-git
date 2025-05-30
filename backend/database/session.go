package database

import (
	"log"
	"time"
)

func StoreSession(sessionID string, userID int, expirationTime time.Time) error {
	// Store session ID in database
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	_, err := db.Exec(`
    INSERT INTO Sessions (user_id, session_token, status, updated_at, expires_at, created_at) 
    VALUES (?, ?, 'active', ?, ?, ?)`,
		userID,
		sessionID,
		currentTime,
		expirationTime.Format("2006-01-02 15:04:05"), // expires_at (correct format)
		currentTime,
	)
	return err
}

func CheckSessionExpiry(userID int) bool {

	var sessionExpiry string
	err := db.QueryRow("SELECT expires_at FROM Sessions WHERE user_id = ?", userID).Scan(&sessionExpiry)
	if err != nil {
		log.Println("No userID found")
		return false
	}
	parsedTime, err := time.Parse("2006-01-02 15:04:05", sessionExpiry)
	if err != nil {
		log.Println("Error parsing session expiry:", err)
		return false
	}

	if parsedTime.After(time.Now()) {
		_, err = db.Exec("UPDATE Sessions SET status = 'expired', updated_at = ? WHERE user_id = ?", time.Now().Format("2006-01-02 15:04:05"), userID)
		if err != nil {
			log.Println("Error updating session expiry:", err)
		}
		return false
	}

	return true
}

// Deleting the session based on sessionID and then any active sessions for that user
func DeleteActiveSession(userID int) error {

	_, err := db.Exec(`
		UPDATE Sessions
		SET status = 'deleted', updated_at = ?
		WHERE user_id = ? AND status = 'active'

	`, time.Now().Format("2006-01-02 15:04:05"), userID)

	return err

}

// GetSessionFromDB retrieves the user ID associated with a valid session cookie
func GetSessionFromDB(sessionID string) (int, error) {
	var userID int
	err := db.QueryRow(`
		SELECT user_id 
		FROM Sessions 
		WHERE session_token = ? 
		AND status = 'active' 
		AND expires_at > datetime('now')`, sessionID).Scan(&userID)
	if err != nil {
		log.Println("No valid session found for the cookie")
		return 0, err
	}
	return userID, nil
}
