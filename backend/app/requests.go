package app

import (
	"log"
	"net/http"
	"social_network/models"
)

func HandleRequests(w http.ResponseWriter, r *http.Request, userID int) {

	var request models.Request
	err := ParseContent(r, &request)
	if err != nil {
		log.Println("Error parsing request:", err)
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Bad Request"})
		return
	}

	log.Println("Parsed request:", request)

	request.Sender.UserID = userID

	if request.Group.GroupID != 0 {
		JoinGroup(w, r, request, userID)
	} else if request.Receiver.UserID > 0 && request.Sender.UserID > 0 {
		// Follow or Unfollow
		if request.Status == "follow" {
			HandleNewFollower(w, r, request)
		} else if request.Status == "unfollow" {
			HandleUnfollow(w, r, request.Sender.UserID, request.Receiver.UserID)
		}
	} else if request.RequestID != 0 && (request.Status == "accepted" || request.Status == "declined") {
		HandleFollowRequest(w, r, request)

	} else {
		log.Println("Error: Invalid request - no group or user IDs provided")
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Bad Request"})
		return
	}

}
