package app

import (
	"encoding/json"
	"log"
	"net/http"
	"social_network/database"
	"social_network/models"
)

func NewPost(w http.ResponseWriter, r *http.Request, userID int) {

	var newPost models.Post
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newPost)
	if err != nil {
		log.Println("error decoding the data")
		ResponseHandler(w, http.StatusBadRequest, "Bad Request")
		return
	}

	if newPost.PostTitle == "" || newPost.PostContent == "" {
		ResponseHandler(w, http.StatusBadRequest, "Title or content cannot be empty")
		return
	}

	err = database.AddPostToDatabase(newPost.PostTitle, newPost.PostContent, userID)
	if err != nil {
		ResponseHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	ResponseHandler(w, http.StatusOK, "Message added to database")

}
