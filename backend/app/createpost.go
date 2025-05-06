package app

import (
	"encoding/json"
	"log"
	"net/http"
	"real-time-forum/database"
	"real-time-forum/models"
)

func FetchCategories(w http.ResponseWriter, r *http.Request) {
	var data []models.CategoryDetails
	var err error
	data, err = database.GetCategories()
	if err != nil {
		log.Println("Error fething categories: ", err)
		ResponseHandler(w, http.StatusInternalServerError, "Internal Server Error")
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func NewPost(w http.ResponseWriter, r *http.Request, userID int) {

	var newPost models.PostDetails
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

	categories := newPost.Categories

	if len(categories) == 0 {
		categories = append(categories, "1") // If no category chosen, give category id 1 (=general)
	}

	var categoryIDs []int
	// Converting categoryIDs to integers and validating that they exists in the database
	for _, cat := range categories {
		categoryID, err := database.HandleCategory(cat)
		if err != nil {
			log.Println("Error handling categoryID in createpost", err)
			ResponseHandler(w, http.StatusBadRequest, "Bad Request")
		}

		categoryIDs = append(categoryIDs, categoryID)
	}

	err = database.AddPostToDatabase(newPost.PostTitle, newPost.PostContent, categoryIDs, userID)
	if err != nil {
		ResponseHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	ResponseHandler(w, http.StatusOK, "Message added to database")

}
