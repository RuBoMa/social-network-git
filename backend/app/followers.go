package app

import (
	"log"
	"net/http"
	"social_network/models"
)

// When someone has pressed "follow"
func HandleNewFollower(w http.ResponseWriter, r *http.Request) {
	followerID := r.URL.Query().Get("follower_id")
	followedID := r.URL.Query().Get("followed_id")
	if followerID == "" || followedID == "" {
		log.Println("Error: follower_id or followed_id not provided")
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Bad Request"})
		return
	}
	// check if profile public or private

	// IF PUBLIC
	// Check if user is already in the follow table
	// Change status / add user to followers table

	// IF PRIVATE
	// Check if user ir already in the follow_requests table
	// Change status / add user to the table as pending

	// Send notification about new follower / follow request

}

// Accept or decline to follow request
func HandleFollowRequest() {

	// IF ACCEPT
	// Change status to accept
	// Change status / add user to Followers table
	// Send notification about acceptance status

	// IF DECLINE
	// Change status to decline

}

// When user pressed "unfollow" on someone
func HandleUnfollow() {

	// Change status on followers table
	// Check if the user is a post_privacy table
	// if yes, change status to deleted

}

func IsProfilePrivate(userID string) (bool, error) {
	var isPrivate bool
	query := `SELECT is_private FROM Users WHERE id = ?`
	err := db.QueryRow(query, userID).Scan(&isPrivate)
	return isPrivate, err
}

func AddFollowRequest(followerID, followedID string) error {
	query := `INSERT INTO Followers (follower_id, followed_id) VALUES (?, ?)`
	_, err := db.Exec(query, followerID, followedID)
	return err
}
func AddFollower(followerID, followedID string) error {
	query := `INSERT INTO Followers (follower_id, followed_id) VALUES (?, ?)`
	_, err := db.Exec(query, followerID, followedID)
	return err
}
func RemoveFollower(followerID, followedID string) error {
	query := `DELETE FROM Followers WHERE follower_id = ? AND followed_id = ?`
	_, err := db.Exec(query, followerID, followedID)
	return err
}
func RemoveFollower(followerID, followedID string) error {
	query := `DELETE FROM Followers WHERE follower_id = ? AND followed_id = ?`
	_, err := db.Exec(query, followerID, followedID)
	return err
}
