package database

import (
	"database/sql"
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
		sessionID, // Using the same UUID for session_token
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

	/* Might be a misstake here..
	_, err = db.Exec("UPDATE Sessions SET status = 'expired', updated_at = ? WHERE user_id = ?", ...)
	according to Chat "You can’t use AND inside SET; you must use commas to separate fields being updated." -TD*/
	if parsedTime.After(time.Now()) {
		_, err = db.Exec("UPDATE Sessions SET status = 'expired' AND updated_at = ? WHERE user_id = ?", time.Now().Format("2006-01-02 15:04:05"), userID)
		if err != nil {
			log.Println("Error updating session expiry:", err)
		}
		return false
	}

	return true
}

// Deleting the session based on sessionID and then any active sessions for that user
func DeleteActiveSession(sessionID string) error {

	var userID int

	// Should it be Sessions or Session? -TD
	err := db.QueryRow(`
		UPDATE Session
		SET status = 'deleted', updated_at = ?
		WHERE id = ? AND status = 'active'
		RETURNING user_id
	`, time.Now().Format("2006-01-02 15:04:05"), sessionID).Scan(&userID)

	// If no active sessions exist
	if err == sql.ErrNoRows {
		return nil
	}

	return err

}

// Changed WHERE to session_token from id. Feels more appropriate since we are storing UUID as session_token
func GetSessionFromDB(sessionID string) (int, error) {
	var userID int
	err := db.QueryRow("SELECT user_id FROM Sessions WHERE session_token = ? AND status = 'active'", sessionID).Scan(&userID)
	if err != nil {
		log.Println("No userID found for the cookie")
		return 0, err
	}
	return userID, nil
}
