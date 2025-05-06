package app

import (
	"log"
	"net/http"
	"real-time-forum/database"
	"time"
)

// Logout logs out the user by deleting the session from the database and setting the session cookie to expire
func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Println("Session cookie not found:", err)
		ResponseHandler(w, http.StatusBadRequest, "No session cookie found")
		return
	}

	err = database.DeleteActiveSession(cookie.Value)
	if err != nil {
		// if return error, it didn't find any sessions to delete
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
	})

	ResponseHandler(w, http.StatusOK, "Logout successful")
}
