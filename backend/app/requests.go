package app

import (
	"net/http"
	"social_network/database"
	"social_network/models"
)

func HandleRequests(w http.ResponseWriter, r *http.Request) {

	var request models.Request
	err := ParseContent(r, &request)
	if err != nil {
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Bad Request"})
		return
	}

	if !ValidIDs(request) {
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Invalid request IDs"})
		return
	}

	if request.Group.GroupID != 0 {
		JoinGroup(w, r, request)
	} else if request.Receiver.UserID > 0 && request.Sender.UserID > 0 {
		// Follow or Unfollow
		HandleNewFollower(w, r, request)
	} else {
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Bad Request"})
		return
	}

}

// ValidRequestIDs checks if the request IDs are valid
// It checks if the senderID, groupID, and receiverID are given and exist in the database
func ValidIDs(request models.Request) bool {
	// If no senderID or it's not valid, return false (senderID is required)
	if request.Sender.UserID == 0 || !database.IsValidUserID(request.Sender.UserID) {
		return false
	}
	// If groupID is provided, check if it's valid
	if request.Group.GroupID != 0 && !database.IsValidGroupID(request.Group.GroupID) {
		return false
	}
	// If receiverID is provided, check if it's valid
	if request.Receiver.UserID != 0 && !database.IsValidUserID(request.Receiver.UserID) {
		return false
	}

	if request.RequestID != 0 && !database.IsValidRequestID(request.RequestID) {
		return false
	}

	return true
}
