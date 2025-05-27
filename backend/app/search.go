package app

import (
	"log"
	"net/http"
	"social_network/database"
	"social_network/models"
)

// Search handles the search functionality for users, groups, posts, and events
// It takes a query string and user ID as parameters and returns the search results
func Search(w http.ResponseWriter, r *http.Request, query string, userID int) {

	var result models.SearchResult
	var err error

	result.Users, err = database.SearchUsers(query, userID)
	if err != nil {
		log.Println("Error searching users:", err)
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal server error"})
		return
	}
	result.Groups, err = database.SearchGroups(query)
	if err != nil {
		log.Println("Error searching groups:", err)
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal server error"})
		return
	}
	result.Posts, err = database.SearchPosts(query, userID)
	if err != nil {
		log.Println("Error searching posts:", err)
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal server error"})
		return
	}
	result.Events, err = database.SearchEvents(query, userID)
	if err != nil {
		log.Println("Error searching events:", err)
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal server error"})
		return
	}

	// Return the search results
	ResponseHandler(w, http.StatusOK, result)
}
