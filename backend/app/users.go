package app

import (
	"net/http"
	"social_network/database"
)

func ServeUsers(w http.ResponseWriter, r *http.Request) {
	users, err := database.GetUsers()
	if err != nil {
		ResponseHandler(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}
	ResponseHandler(w, http.StatusOK, users)
}
