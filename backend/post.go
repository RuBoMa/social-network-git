package backend

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// PostHandler handles requests to view a specific post
func PostPage(w http.ResponseWriter, r *http.Request, userID int) {

	pathParts := strings.Split(r.URL.Path, "/")
	var postID int
	var err error
	var addedItem string

	// Remove empty strings from pathParts
	var filteredParts []string
	for _, part := range pathParts {
		if part != "" {
			filteredParts = append(filteredParts, part)
		}
	}
	pathParts = filteredParts
	// Structure should be /api/post/post_id/comment
	if len(pathParts) >= 3 && pathParts[1] == "post" {

		postID, err = strconv.Atoi(pathParts[2])
		if err != nil {
			log.Println("Error converting postID to int:", err)
			ResponseHandler(w, http.StatusNotFound, "Page Not Found")
			return
		}

		// Check if comment or vote is added
		if len(pathParts) >= 4 {
			addedItem = pathParts[3]
		}
	}

	valid := ValidatePostID(postID)
	if !valid {
		log.Println("Invalid postID: ", postID)
		ResponseHandler(w, http.StatusNotFound, "Page Not Found")
		return
	}

	switch r.Method {
	case http.MethodGet:
		HandlePostPageGet(w, r, postID, userID)
	case http.MethodPost:
		if addedItem == "comment" {
			HandleComment(w, r, postID, userID)
		} else if addedItem == "vote" {
			HandleVote(w, r, postID, userID)
		} else {
			ResponseHandler(w, http.StatusBadRequest, "Bad Request")
		}
	default:
		ResponseHandler(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}
}

// HandlePostPageGet handles get requests to the post page
func HandlePostPageGet(w http.ResponseWriter, r *http.Request, postID, userID int) {
	post, err := GetPostDetails(postID, userID)
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

	var newComment CommentDetails
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newComment)
	if err != nil {
		log.Println("error decoding the data")
		ResponseHandler(w, http.StatusBadRequest, "Bad Request")
		return
	}

	if newComment.Content != "" {
		// Insert comment into the database
		err := AddComment(postID, newComment.Content, userID)
		if err != nil {
			ResponseHandler(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	HandlePostPageGet(w, r, postID, userID)
}

// HandlePostPagePost handles post requests to the post page
func HandleVote(w http.ResponseWriter, r *http.Request, postID, userID int) {

	var newVote VoteDetails
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
		exists := ValidateCommentID(newVote.CommentID)
		if !exists {
			log.Println("CommentID doesn't exist", comment_id)
			ResponseHandler(w, http.StatusBadRequest, "Bad Request")
			return
		}
		post_id = 0
	}

	err = AddVotes(userID, post_id, comment_id, likeType)
	if err != nil {
		log.Printf("Error adding votes to the database: userID %d, postID %d, commentID %d, like type %d\n", userID, post_id, comment_id, likeType)
		ResponseHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	HandlePostPageGet(w, r, postID, userID)
}

// ValidatePostID checks if a post with the given ID exists in the database
func ValidatePostID(postID int) bool {
	var post int
	err := db.QueryRow("SELECT id FROM Post WHERE id = ?", postID).Scan(&post)
	if err != nil {
		log.Println("Error scanning postID:", err)
		return false
	}
	return true
}

// ValidateCommentID checks if a comment with the given ID exists in the database
func ValidateCommentID(commentID int) bool {
	var comment int
	err := db.QueryRow("SELECT id FROM Comment WHERE id = ?", commentID).Scan(&comment)
	if err != nil {
		log.Println("Error scanning commentID:", err)
		return false
	}
	return true
}
