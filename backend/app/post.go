package app

import (
	"encoding/json"
	"log"
	"net/http"
	"social_network/database"
	"social_network/models"
)

// HandlePostPageGet handles get requests to the post page
func HandlePostPageGet(w http.ResponseWriter, r *http.Request, postID, userID int) {
	post, err := database.GetPostDetails(postID, userID)
	if err != nil {
		log.Println("Error fetching post details:", err)
		ResponseHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(post); err != nil {
		ResponseHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
}

// HandlePostPagePost handles post requests to the post page
func HandleComment(w http.ResponseWriter, r *http.Request, postID, userID int) {

	var newComment models.CommentDetails
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newComment)
	if err != nil {
		log.Println("error decoding the data")
		ResponseHandler(w, http.StatusBadRequest, "Bad Request")
		return
	}

	if newComment.CommentContent != "" {
		// Insert comment into the database
		err := database.AddComment(postID, newComment.CommentContent, userID)
		if err != nil {
			ResponseHandler(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	HandlePostPageGet(w, r, postID, userID)
}
