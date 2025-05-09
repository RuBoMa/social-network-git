package server

import (
	"encoding/json"
	"net/http"
	"social_network/app"
)

// HomePage handles the rendering of the home page
func HandleFeed(w http.ResponseWriter, r *http.Request, userID int) {

	posts, err := app.FetchFeed(userID)
	if err != nil {
		app.ResponseHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	// Return posts as JSON
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		app.ResponseHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

}
