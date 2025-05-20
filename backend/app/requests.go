package app

import (
	"log"
	"net/http"
	"social_network/database"
	"social_network/models"
)

func HandleRequests(w http.ResponseWriter, r *http.Request) {

	var request models.Request
	err := ParseContent(r, &request)
	if err != nil {
		log.Println("Error parsing request:", err)
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Bad Request"})
		return
	}
	log.Println("Request:", request)

	if !ValidIDs(request) {
		log.Println("Error: Invalid request IDs")
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Invalid request IDs"})
		return
	}

	if request.Group.GroupID != 0 {
		JoinGroup(w, r, request)
	} else if request.Receiver.UserID > 0 && request.Sender.UserID > 0 {
		// Follow or Unfollow
		if request.Status == "follow" {
			HandleNewFollower(w, r, request)
		} else if request.Status == "unfollow" {
			HandleUnfollow(w, r, request.Sender.UserID, request.Receiver.UserID)
		}
	} else {
		log.Println("Error: Invalid request - no group or user IDs provided")
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
