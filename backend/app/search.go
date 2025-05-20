package app

import (
	"net/http"
	"social_network/database"
	"social_network/models"
)

func Search(w http.ResponseWriter, r *http.Request, query string, userID int) {

	var result models.SearchResult
	var err error

	result.Users, err = database.SearchUsers(query)
	result.Groups, err = database.SearchGroups(query)
	result.Posts, err = database.SearchPosts(query, userID)
	result.Events, err = database.SearchEvents(query, userID)

	if err != nil {
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Database error"})
		return
	}

	// Return the search results
	ResponseHandler(w, http.StatusOK, result)
}
