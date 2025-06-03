package app

import (
	"log"
	"net/http"
	"social_network/database"
)

// HandleFeed processes the request to fetch the user's or group's feed.
func HandleFeed(w http.ResponseWriter, r *http.Request, userID, groupID int) {

	if groupID != 0 {
		userID = 0
	}

	posts, err := database.GetPosts(userID, groupID)
	if err != nil {
		log.Println("Error fetching feed:", err)
		ResponseHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	ResponseHandler(w, http.StatusOK, posts)

}
