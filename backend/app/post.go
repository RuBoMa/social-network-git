package app

import (
	"encoding/json"
	"log"
	"net/http"
	"real-time-forum/backend/database"
	"real-time-forum/backend/models"
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

	if newComment.Content != "" {
		// Insert comment into the database
		err := database.AddComment(postID, newComment.Content, userID)
		if err != nil {
			ResponseHandler(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	HandlePostPageGet(w, r, postID, userID)
}

// HandlePostPagePost handles post requests to the post page
func HandleVote(w http.ResponseWriter, r *http.Request, postID, userID int) {

	var newVote models.VoteDetails
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newVote)
	if err != nil {
		log.Println("error decoding the data")
		ResponseHandler(w, http.StatusBadRequest, "Bad Request")
		return
	}
	// Insert vote into the database
	var likeType int
	var post_id int
	var comment_id int
	if newVote.Vote == "like" {
		likeType = 1
	} else if newVote.Vote == "dislike" {
		likeType = 2
	} else {
		log.Println("Invalid vote value: ", newVote.Vote)
		ResponseHandler(w, http.StatusBadRequest, "Bad Request")
		return
	}
	// Check if the vote is for a post or a comment
	if newVote.PostID != 0 {
		comment_id = 0
		post_id = newVote.PostID
	} else {
		exists := database.ValidateCommentID(newVote.CommentID)
		if !exists {
			log.Println("CommentID doesn't exist", comment_id)
			ResponseHandler(w, http.StatusBadRequest, "Bad Request")
			return
		}
		post_id = 0
	}

	err = database.AddVotes(userID, post_id, comment_id, likeType)
	if err != nil {
		log.Printf("Error adding votes to the database: userID %d, postID %d, commentID %d, like type %d\n", userID, post_id, comment_id, likeType)
		ResponseHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	HandlePostPageGet(w, r, postID, userID)
}
