package app

import (
	"database/sql"
	"log"
	"net/http"
	"social_network/database"
	"social_network/models"
)

var db *sql.DB

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
		err = database.AddNotificationIntoDB("follow_request", request, models.Event{})
		if err != nil {
			log.Println("Error saving notification:", err)
			// Currently not crashing the server if notification fails
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
	followedID := request.Sender.UserID
	followerID := request.Receiver.UserID
	action := request.Status

	if followerID == 0 || followedID == 0 || action == "" {
		log.Println("Error: follower_id, followed_id or status not provided")
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Bad Request"})
		return
	}

	err := database.UpdateRequestStatus(request)
	if err != nil {
		log.Println("Error updating follow request status:", err)
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
		return
	}

	if action == "accept" {
		err = database.AddFollower(followerID, followedID)
		if err != nil {
			log.Println("Error adding follower:", err)
			ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
			return
		}
		ResponseHandler(w, http.StatusOK, models.Response{Message: "Follow request accepted"})
	} else if action != "decline" {
		log.Println("Error: invalid action")
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Bad Request"})
		return
	}

	ResponseHandler(w, http.StatusOK, models.Response{Message: "Follow request declined"})
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
