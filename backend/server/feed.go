package server

import (
	"encoding/json"
	"net/http"
	"real-time-forum/backend/app"
)

// HomePage handles the rendering of the home page
func HandleFeed(w http.ResponseWriter, r *http.Request) {

	posts, err := app.FetchFeed(0)
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
