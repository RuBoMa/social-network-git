package app

import (
	"log"
	"net/http"
	"social_network/database"
	"social_network/models"
)

// HandleNewFollower handles the logic when a user follows another user
func HandleNewFollower(w http.ResponseWriter, r *http.Request, request models.Request) {
	followerID := request.Sender.UserID
	followedID := request.Receiver.UserID

	if followerID == 0 || followedID == 0 {
		log.Println("Error: follower_id or followed_id not provided")
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Bad Request"})
		return
	}
	// Check if the profile is public or private
	isPrivate, err := database.IsProfilePrivate(followedID)
	if err != nil {
		log.Println("Error checking profile privacy:", err)
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
		return
	}

	if isPrivate {
		request.RequestID, err = database.AddRequestIntoDB(request)
		if err != nil {
			log.Println("Error adding follow request:", err)
			ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
			return
		}
		notificationIDs, err := database.AddNotificationIntoDB(models.NotifFollowRequest, request, models.Event{})
		if err != nil {
			log.Println("Error saving notification:", err)
			// Currently not crashing the server if notification fails
		}
		for _, notificationID := range notificationIDs {
			ServeNotification(notificationID)
		}
		ResponseHandler(w, http.StatusOK, models.Response{Message: "Follow request sent"})

	} else {

		err = database.AddFollower(followerID, followedID)
		if err != nil {
			log.Println("Error adding follower:", err)
			ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
			return
		}
		ResponseHandler(w, http.StatusOK, models.Response{Message: "Followed successfully"})
	}
}

// HandleFollowRequest handles the logic when a user accepts or declines a follow request
func HandleFollowRequest(w http.ResponseWriter, r *http.Request, request models.Request) {

	if request.Status == "" {
		log.Println("Error: follower_id, followed_id or status not provided")
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Bad Request"})
		return
	}

	isValid := database.IsValidRequestID(request.RequestID)
	if !isValid {
		log.Println("Invalid request ID:", request.RequestID)
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Invalid Request ID"})
		return
	}

	err := database.UpdateRequestStatus(request)
	if err != nil {
		log.Println("Error updating follow request status:", err)
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
		return
	}

	if request.Status == "accepted" {
		request, err = database.GetRequestByID(request.RequestID)
		if err != nil {
			log.Println("Error fetching request by ID:", err)
			ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
			return
		}
		err = database.AddFollower(request.Sender.UserID, request.Receiver.UserID)
		if err != nil {
			log.Println("Error adding follower:", err)
			ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
			return
		}
	}
	ResponseHandler(w, http.StatusOK, models.Response{Message: "Follow request updated successfully"})
}

// HandleUnfollow handles the logic when a user unfollows another user
func HandleUnfollow(w http.ResponseWriter, r *http.Request, followerID, followedID int) {

	if followerID == 0 || followedID == 0 {
		log.Println("Error: follower_id or followed_id not provided")
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Bad Request"})
		return
	}
	err := database.RemoveFollower(followerID, followedID)
	if err != nil {
		log.Println("Error removing follower:", err)
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
		return
	}
	err = database.RemoveFromPostPrivacy(followerID, followedID)
	if err != nil {
		log.Println("Error removing from post privacy:", err)
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
		return
	}
	ResponseHandler(w, http.StatusOK, models.Response{Message: "Unfollowed successfully"})
}

// GetFollowers handles the logic to get the followers of a user
func GetFollowers(w http.ResponseWriter, userID int) {

	followers, err := database.GetFollowers(userID)
	if err != nil {
		log.Println("Error fetching followers:", err)
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
		return
	}

	ResponseHandler(w, http.StatusOK, followers)
}

// GetFollowing handles the logic to get the users that a user is following
func GetFollowing(w http.ResponseWriter, userID int) {
	following, err := database.GetFollowing(userID)
	if err != nil {
		log.Println("Error fetching following:", err)
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
		return
	}

	ResponseHandler(w, http.StatusOK, following)
}
