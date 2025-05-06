package database

import (
	"log"
	"time"
)

func StoreSession(sessionID string, userID int, expirationTime time.Time) error {
	// Store session ID in database
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	_, err := db.Exec(`
    INSERT INTO Session (id, user_id, created_at, updated_at, expires_at, last_access) 
    VALUES (?, ?, ?, ?, ?, ?)`,
		sessionID, // Using the same UUID for session_token
		userID,
		currentTime,
		currentTime,
		expirationTime.Format("2006-01-02 15:04:05"), // expires_at (correct format)
		currentTime,
	)
	return err
}

func CheckSessionExpiry(userID int) bool {

	var sessionExpiry string
	err := db.QueryRow("SELECT expires_at FROM Session WHERE user_id = ?", userID).Scan(&sessionExpiry)
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
		_, err = db.Exec("UPDATE Session SET status = 'expired' AND updated_at = ? WHERE user_id = ?", time.Now().Format("2006-01-02 15:04:05"), userID)
		if err != nil {
			log.Println("Error updating session expiry:", err)
		}
		return false
	}

	return true
}

func RefreshLastAccess(userID int) {
	lastAccessed := time.Now()

	_, err := db.Exec("UPDATE Session SET last_access = ? WHERE user_id = ? AND status = 'active'", lastAccessed.Format("2006-01-02 15:04:05"), userID)
	if err != nil {
		log.Println("Error updating session expiry:", err)
		return
	}
}

// Deleting the session based on sessionID and then any active sessions for that user
func DeleteActiveSession(sessionID string) error {

	var userID int
	err := db.QueryRow(`
	UPDATE Session 
	SET status = 'deleted', updated_at = ?
	WHERE id = ? AND status = 'active'
	RETURNING user_id
	`, time.Now().Format("2006-01-02 15:04:05"), sessionID).Scan(&userID)

	if err == nil {
		_, err = db.Exec(`
		UPDATE Session
		SET status = 'deleted', updated_at = ?
		WHERE user_id = ? AND status = 'active'
	`, time.Now().Format("2006-01-02 15:04:05"), userID)
	}

	return err

}

func GetSessionFromDB(sessionID string) (int, error) {
	var userID int
	err := db.QueryRow("SELECT user_id FROM Session WHERE id = ? AND status = 'active'", sessionID).Scan(&userID)
	if err != nil {
		log.Println("No userID found for the cookie")
		return 0, err
	}
	return userID, nil
}

// updateSessionExpiry updates session expiry and cookie expiry based on the last access
func UpdateSessionExpiry(userID int) (time.Time, error) {
	var lastAccessed string
	var expiresAt string

	err := db.QueryRow("SELECT last_access, expires_at FROM Session WHERE user_id = ? AND status = 'active'", userID).Scan(&lastAccessed, &expiresAt)
	if err != nil {
		log.Println("No userID found:", err)
		return time.Time{}, err
	}

	parsedLastAccess, err := time.Parse("2006-01-02 15:04:05", lastAccessed)
	if err != nil {
		log.Println("Error parsing last_access:", err)
		return time.Time{}, err
	}

	parsedExpiresAt, err := time.Parse("2006-01-02 15:04:05", expiresAt)
	if err != nil {
		log.Println("Error parsing expires_at:", err)
		return time.Time{}, err
	}

	// If last_access and expires_at are not the same, update expires_at and session
	if !parsedLastAccess.Equal(parsedExpiresAt) {
		_, err := db.Exec("UPDATE Session SET expires_at = ? WHERE user_id = ?", lastAccessed, userID)
		if err != nil {
			log.Println("Error updating session expiry:", err)
			return time.Time{}, err
		}

	}
	return parsedLastAccess, nil
}
