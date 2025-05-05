package backend

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func Authenticate(w http.ResponseWriter, loggedIn bool, userID int) {
	status := http.StatusUnauthorized
	message := "No current sessions"

	if loggedIn {
		refreshLastAccess(userID)
		status = http.StatusOK
		message = strconv.Itoa(userID)
	}

	ResponseHandler(w, status, message)
}

// CreateSession creates a new session for the user and stores it in the database
func CreateSession(w http.ResponseWriter, userID int) error {

	if userID == 0 {
		return fmt.Errorf("userID is 0")
	}
	//First check for and delete any existing sessions for this user
	_, err := db.Exec("UPDATE Session SET status = 'deleted', updated_at = ? WHERE user_id = ? AND status = 'active'",
		time.Now().Format("2006-01-02 15:04:05"), userID)
	if err != nil {
		return err
	}
	sessionID := uuid.NewString()
	expirationTime := time.Now().Add(30 * time.Minute)
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  expirationTime,
		HttpOnly: true, // Prevent JavaScript from accessing the cookie
		Path:     "/",
	})

	// Store session ID in database
	_, err = db.Exec(`
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

// VerifySession checks if the session ID exists in the database
func VerifySession(r *http.Request) (bool, int) {
	var userID int
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return false, 0
	}

	err = db.QueryRow("SELECT user_id FROM Session WHERE id = ? AND status = 'active'", cookie.Value).Scan(&userID)
	if err != nil {
		log.Println("No userID found for the cookie")
		return false, 0
	}

	return true, userID
}

func refreshLastAccess(userID int) {
	lastAccessed := time.Now()

	_, err := db.Exec("UPDATE Session SET last_access = ? WHERE user_id = ? AND status = 'active'", lastAccessed.Format("2006-01-02 15:04:05"), userID)
	if err != nil {
		log.Println("Error updating session expiry:", err)
		return
	}
}

// updateSessionExpiry updates session expiry and cookie expiry based on the last access
func updateSessionExpiry(userID int, w http.ResponseWriter) {
	var lastAccessed string
	var expiresAt string

	err := db.QueryRow("SELECT last_access, expires_at FROM Session WHERE user_id = ? AND status = 'active'", userID).Scan(&lastAccessed, &expiresAt)
	if err != nil {
		log.Println("No userID found:", err)
		return
	}

	parsedLastAccess, err := time.Parse("2006-01-02 15:04:05", lastAccessed)
	if err != nil {
		log.Println("Error parsing last_access:", err)
		return
	}

	parsedExpiresAt, err := time.Parse("2006-01-02 15:04:05", expiresAt)
	if err != nil {
		log.Println("Error parsing expires_at:", err)
		return
	}

	// If last_access and expires_at are not the same, update expires_at and session
	if !parsedLastAccess.Equal(parsedExpiresAt) {
		_, err := db.Exec("UPDATE Session SET expires_at = ? WHERE user_id = ?", lastAccessed, userID)
		if err != nil {
			log.Println("Error updating session expiry:", err)
			return
		}

		cookieExpiry := parsedLastAccess.Add(30 * time.Minute)
		http.SetCookie(w, &http.Cookie{
			Name:    "session_expiry",
			Value:   cookieExpiry.Format("2006-01-02 15:04:05"),
			Expires: cookieExpiry,
			Path:    "/",
		})
	}
}

func checkSessionExpiry(userID int) bool {

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

// Handler to verify or expire session
func SessionHandler(w http.ResponseWriter, loggedIn bool, userID int) {
	status := http.StatusOK
	message := "Session active"

	if !loggedIn {
		status = http.StatusUnauthorized
		message = "Session expired"
	} else {
		updateSessionExpiry(userID, w)
		activeSession := checkSessionExpiry(userID)
		if !activeSession {
			status = http.StatusUnauthorized
			message = "Session expired"
		}
	}
	ResponseHandler(w, status, message)
}
