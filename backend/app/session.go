package app

import (
	"fmt"
	"net/http"
	"social_network/database"
	"time"

	"github.com/google/uuid"
)

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
