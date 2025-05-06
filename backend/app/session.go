package app

import (
	"fmt"
	"log"
	"net/http"
	"real-time-forum/database"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func Authenticate(w http.ResponseWriter, loggedIn bool, userID int) {
	status := http.StatusUnauthorized
	message := "No current sessions"

	if loggedIn {
		database.RefreshLastAccess(userID)
		status = http.StatusOK
		message = strconv.Itoa(userID)
	}

	ResponseHandler(w, status, message)
}

// CreateSession creates a new session for the user and stores it in the database
func CreateSession(w http.ResponseWriter, r *http.Request, userID int) error {

	if userID == 0 {
		return fmt.Errorf("userID is 0")
	}

	cookie, err := r.Cookie("session_id")
	if err == nil {
		err := database.DeleteActiveSession(cookie.Value)
		if err != nil {
			return err
		}
	}

	sessionID := uuid.NewString()
	expirationTime := time.Now().Add(30 * time.Minute)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  expirationTime,
		HttpOnly: true, // Prevent JavaScript from accessing the cookie
		Path:     "/",
	})

	err = database.StoreSession(sessionID, userID, expirationTime)

	return err

}

// VerifySession checks if the session ID exists in the database
func VerifySession(r *http.Request) (bool, int) {

	cookie, err := r.Cookie("session_id")
	if err != nil {
		return false, 0
	}

	userID, err := database.GetSessionFromDB(cookie.Value)
	if err != nil {
		return false, 0
	}

	return true, userID
}

// updateSessionExpiry updates session expiry and cookie expiry based on the last access
func UpdateSessionExpiry(userID int, w http.ResponseWriter) {

	parsedLastAccess, err := database.UpdateSessionExpiry(userID)
	if err != nil {
		log.Println("Error updating session expiry on database: ", err)
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

// Handler to verify or expire session
func SessionHandler(w http.ResponseWriter, loggedIn bool, userID int) {
	status := http.StatusOK
	message := "Session active"

	if !loggedIn {
		status = http.StatusUnauthorized
		message = "Session expired"
	} else {
		UpdateSessionExpiry(userID, w)
		activeSession := database.CheckSessionExpiry(userID)
		if !activeSession {
			status = http.StatusUnauthorized
			message = "Session expired"
		}
	}
	ResponseHandler(w, status, message)
}
