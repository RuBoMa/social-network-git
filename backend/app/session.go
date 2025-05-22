package app

import (
	"encoding/json"
	"fmt"
	"log"
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
		MaxAge:   1800,
		HttpOnly: true, // Prevent JavaScript from accessing the cookie
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
	})

	err = database.StoreSession(sessionID, userID, expirationTime)

	return err

}

// VerifySession checks if the session ID exists in the database
func VerifySession(r *http.Request) (bool, int) {
	token := r.URL.Query().Get("token")
	if token == "" {

		cookie, err := r.Cookie("session_id")
		if err != nil {
			log.Println("Session cookie not found:", err)
			return false, 0
		}
		token = cookie.Value
	}

	userID, err := database.GetSessionFromDB(token)
	if err != nil {
		log.Println("Error getting session from DB:", err)
		return false, 0
	}
	log.Println("Session verified for user ID:", userID)
	return true, userID
}
func VerifySessionToken(token string) (bool, int) {
	if token == "" {
		return false, 0
	}

	userID, err := database.GetSessionFromDB(token)
	if err != nil {
		log.Println("Error getting session from DB:", err)
		return false, 0
	}
	log.Println("Session verified for user ID:", userID)
	return true, userID
}

func GetSessionHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := database.GetSessionFromDB(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // adjust to your frontend URL
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":  cookie.Value,
		"userId": userID,
	})
}
