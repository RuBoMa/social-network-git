package app

import (
	"fmt"
	"log"
	"net/http"
	"social_network/database"
	"time"

	"github.com/google/uuid"
)
// CreateSession creates a new session for the user and stores it in the database
func CreateSession(w http.ResponseWriter, r *http.Request, userID int) (string, error) {

	if userID == 0 {
		return "", fmt.Errorf("userID is 0")
	}

	err := database.DeleteActiveSession(userID)
	if err != nil {
		log.Println("Error deleting active session:", err)
		return "", err
	}

	sessionID := uuid.NewString()
	expirationTime := time.Now().Add(24 * time.Hour)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  expirationTime,
		MaxAge:   24 * 60 * 60, // 1 day
		HttpOnly: true,         // Prevent JavaScript from accessing the cookie
		Path:     "/",
	})

	err = database.StoreSession(sessionID, userID, expirationTime)

	return sessionID, err

}

// VerifySession checks if the session ID exists in the database
func VerifySession(r *http.Request) (bool, int) {

	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Println("Session cookie not found:", err)
		return false, 0
	}

	userID, err := database.GetSessionFromDB(cookie.Value)
	if err != nil {
		log.Println("Error getting session from DB:", err)
		return false, 0
	}
	log.Println("Session verified for user ID:", userID)
	return true, userID
}

func VerifySessionHandler(w http.ResponseWriter, r *http.Request) {
	ok, _ := VerifySession(r)
	if !ok {
		log.Println("Session verification failed in session handler")
		ResponseHandler(w, http.StatusUnauthorized, "Session not found or expired")
		return
	}
	log.Println("Session verified successfully in session handler")
	ResponseHandler(w, http.StatusOK, "Session verified successfully")
}
