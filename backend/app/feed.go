package app

import (
	"net/http"
)

// HomePage handles the rendering of the home page
func HandleFeed(w http.ResponseWriter, r *http.Request, userID int) {

	posts, err := FetchFeed(userID)
	if err != nil {
		ResponseHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	ResponseHandler(w, http.StatusOK, posts)

}
