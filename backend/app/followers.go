package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"social_network/models"
	"strings"
)

var db *sql.DB

// HandleNewFollower handles the logic when a user follows another user
func HandleNewFollower(w http.ResponseWriter, r *http.Request) {
	followerID := r.URL.Query().Get("follower_id")
	followedID := r.URL.Query().Get("followed_id")

	if followerID == "" || followedID == "" {
		log.Println("Error: follower_id or followed_id not provided")
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Bad Request"})
		return
	}
	// Check if the profile is public or private
	isPrivate, err := IsProfilePrivate(followedID)
	if err != nil {
		log.Println("Error checking profile privacy:", err)
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
		return
	}

	if isPrivate {
		err = AddFollowRequest(followerID, followedID)
		if err != nil {
			log.Println("Error adding follow request:", err)
			ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
			return
		}
		go SendFollowRequestNotification(followerID, followedID)
		ResponseHandler(w, http.StatusOK, models.Response{Message: "Follow request sent"})
	} else {
		err = AddFollower(followerID, followedID)
		if err != nil {
			log.Println("Error adding follower:", err)
			ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
			return
		}
		go SendNewFollowerNotification(followerID, followedID)
		ResponseHandler(w, http.StatusOK, models.Response{Message: "Followed successfully"})
	}
}

// HandleFollowRequest handles the logic when a user accepts or declines a follow request
func HandleFollowRequest(w http.ResponseWriter, r *http.Request) {
	followerID := r.URL.Query().Get("follower_id")
	followedID := r.URL.Query().Get("followed_id")
	action := r.URL.Query().Get("action") // accept or decline

	if followerID == "" || followedID == "" || action == "" {
		log.Println("Error: follower_id, followed_id or action not provided")
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Bad Request"})
		return
	}

	if action == "accept" {
		err := AcceptFollowRequest(followerID, followedID)
		if err != nil {
			log.Println("Error accepting follow request:", err)
			ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
			return
		}
		err = AddFollower(followerID, followedID)
		if err != nil {
			log.Println("Error adding follower:", err)
			ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
			return
		}
		ResponseHandler(w, http.StatusOK, models.Response{Message: "Follow request accepted"})
	} else if action == "decline" {
		err := DeclineFollowRequest(followerID, followedID)
		if err != nil {
			log.Println("Error declining follow request:", err)
			ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
			return
		}
		ResponseHandler(w, http.StatusOK, models.Response{Message: "Follow request declined"})
	} else {
		log.Println("Error: invalid action")
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Bad Request"})
		return
	}
}

// HandleUnfollow handles the logic when a user unfollows another user
func HandleUnfollow(w http.ResponseWriter, r *http.Request) {
	followerID := r.URL.Query().Get("follower_id")
	followedID := r.URL.Query().Get("followed_id")
	if followerID == "" || followedID == "" {
		log.Println("Error: follower_id or followed_id not provided")
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Bad Request"})
		return
	}
	err := RemoveFollower(followerID, followedID)
	if err != nil {
		log.Println("Error removing follower:", err)
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
		return
	}
	err = RemoveFromPostPrivacy(followerID, followedID)
	if err != nil {
		log.Println("Error removing from post privacy:", err)
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
		return
	}
}

// Check if the profile is public or private
func IsProfilePrivate(userID string) (bool, error) {
	var isPublic bool
	query := `SELECT is_public FROM Users WHERE id = ?`
	err := db.QueryRow(query, userID).Scan(&isPublic)
	return !isPublic, err
}
// AddFollowRequest adds a follow request to the database
func AddFollowRequest(followerID, followedID string) error {
	query := `INSERT INTO Followers (follower_id, followed_id, status) VALUES (?, ?, 'pending')`
	_, err := db.Exec(query, followerID, followedID)
	return err
}
// AddFollower adds a follower to the database
func AddFollower(followerID, followedID string) error {
	query := `INSERT INTO Followers (follower_id, followed_id) VALUES (?, ?)`
	_, err := db.Exec(query, followerID, followedID)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return fmt.Errorf("user is already a following")
		}
		return err
	}
	return nil
}
// AcceptFollowRequest accepts a follow request and updates the status in the database
func AcceptFollowRequest(followerID, followedID string) error {
	query := `UPDATE Followers SET status = 'accepted' WHERE follower_id = ? AND followed_id = ?`
	_, err := db.Exec(query, followerID, followedID)
	return err
}
// DeclineFollowRequest declines a follow request and updates the status in the database
func DeclineFollowRequest(followerID, followedID string) error {
	query := `UPDATE Followers SET status = 'declined' WHERE follower_id = ? AND followed_id = ?`
	_, err := db.Exec(query, followerID, followedID)
	return err
}
// RemoveFollower removes a follower from the database
func RemoveFollower(followerID, followedID string) error {
	query := `DELETE FROM Followers WHERE follower_id = ? AND followed_id = ?`
	_, err := db.Exec(query, followerID, followedID)
	return err
}
// RemoveFromPostPrivacy removes a user from the post privacy list
func RemoveFromPostPrivacy(followerID, followedID string) error {
	query := `DELETE FROM PostPrivacy WHERE user_id = ? AND post_id = ?`
	_, err := db.Exec(query, followedID, followerID)
	return err
}

