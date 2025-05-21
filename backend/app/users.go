package app

import (
	"net/http"
	"social_network/database"
)

// ServeUsers fetches all users from the database
func ServeUsers(w http.ResponseWriter, r *http.Request) {
	users, err := database.GetUsers()
	if err != nil {
		ResponseHandler(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}
	ResponseHandler(w, http.StatusOK, users)
}

// ServeNonGroupMembers fetches users who are not members of a specific group and have not outstanding invitations/requests
func ServeNonGroupMembers(w http.ResponseWriter, r *http.Request, groupID int) {
	users, err := database.GetNonGroupMembers(groupID)
	if err != nil {
		ResponseHandler(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}
	ResponseHandler(w, http.StatusOK, users)
}
