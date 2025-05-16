package app

import (
	"log"
	"net/http"
)

// HomePage handles the rendering of the home page
func HandleFeed(w http.ResponseWriter, r *http.Request, userID, groupID int) {

	log.Println("Fetching feed for userID:", userID, "and groupID:", groupID)
	if groupID != 0 {
		userID = 0
	}
	log.Println("Fetching feed for userID:", userID, "and groupID:", groupID)

	posts, err := FetchFeed(userID, groupID)
	if err != nil {
		ResponseHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	ResponseHandler(w, http.StatusOK, posts)

}
