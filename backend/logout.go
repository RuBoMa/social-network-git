package backend

import (
	"log"
	"net/http"
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

	var userID int
	err = db.QueryRow("SELECT user_id FROM Session WHERE id = ?", cookie.Value).Scan(&userID)
	if err != nil {
		log.Println("Error fetching user_id from session:", err)
		ResponseHandler(w, http.StatusInternalServerError, "Failed to retrieve session data")
		return
	}

	_, err = db.Exec("UPDATE Session SET status = 'deleted', updated_at = ? WHERE id = ? AND status = 'active'",
		time.Now().Format("2006-01-02 15:04:05"), cookie.Value)
	if err != nil {
		log.Println("Error deleting session:", err)
		ResponseHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
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
