package app

import (
	"net/http"
	"social_network/models"
)

func HandleRequests(w http.ResponseWriter, r *http.Request) {

	var request models.Request
	err := ParseContent(r, &request)
	if err != nil {
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Bad Request"})
		return
	}

	if request.Group.GroupID != 0 {
		JoinGroup(w, r, request)
	} else if request.Receiver.UserID > 0 && request.Sender.UserID > 0 {
		// Follow or Unfollow
	} else {
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Bad Request"})
		return
	}



}
