package app

import (
	"encoding/json"
	"log"
	"net/http"
	"social_network/database"
	"social_network/models"
	"strconv"
)

// FetchFeed retrieves the posts for a given user or group.
// func FetchFeed(userID, groupID int) ([]models.Post, error) {
// 	posts, err := database.GetPosts(userID, groupID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return posts, nil
// }

// HandlePostGet handles get requests to a specific post
// It retrieves the post details and comments associated with it
func HandlePostGet(w http.ResponseWriter, r *http.Request, postID, userID int) {

	if !database.ValidatePostID(postID) {
		log.Println("Invalid postID: ", postID)
		ResponseHandler(w, http.StatusNotFound, models.Response{Message: "Page Not Found"})
		return
	}

	canView := database.CheckPostPrivacy(postID, userID)
	if !canView {
		log.Println("User does not have permission to view this post")
		ResponseHandler(w, http.StatusForbidden, models.Response{Message: "You do not have permission to view this post"})
		return
	}
	post, err := database.GetPostDetails(postID)
	if err != nil {
		log.Println("Error fetching post details:", err)
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
		return
	}

	// Log successful retrieval
	log.Printf("Successfully retrieved post %d for user %d", postID, userID)
	ResponseHandler(w, http.StatusOK, post)
}

// NewPost handles post requests to create a new post
// It expects a multipart form with the post title, content, privacy setting, and an optional image
// It also associates the post with a group if provided
func NewPost(w http.ResponseWriter, r *http.Request, userID int) {

	var newPost models.Post

	err := ParseContent(r, &newPost)
	if err != nil {
		log.Println("Error parsing the new post data:", err)
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Bad Request"})
		return
	}

	// If multipart form, get the form values
	if newPost.PostTitle == "" {
		newPost.PostTitle = r.FormValue("post_title")
		newPost.PostContent = r.FormValue("post_content")
		newPost.Privacy = r.FormValue("privacy")
		customUsersStr := r.FormValue("custom_users")
		if customUsersStr != "" {

			var userIDs []int
			err := json.Unmarshal([]byte(customUsersStr), &userIDs)
			if err != nil {
				log.Println("Failed to parse user IDs JSON:", err)
				ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Invalid user IDs"})
				return
			}

			for _, userID := range userIDs {
				newPost.CustomUsers = append(newPost.CustomUsers, models.User{UserID: userID})
			}

			log.Println("Custom users: ", newPost.CustomUsers)
		}
		groupIDSstr := r.FormValue("group_id")
		newPost.PostImage = SaveUploadedFile(r, "post_image", "post")

		if groupIDSstr != "" {
			groupID, err := strconv.Atoi(groupIDSstr)
			if err != nil {
				ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Invalid group ID"})
				return
			}
			newPost.Group.GroupID = groupID
		}
	}

	if newPost.PostTitle == "" || newPost.PostContent == "" || newPost.Privacy == "" {
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Title, content or privacy setting cannot be empty"})
		return
	}

	err = database.AddPostIntoDB(newPost, userID)
	if err != nil {
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
		return
	}

	ResponseHandler(w, http.StatusOK, models.Response{Message: "Message added to database"})

}

// NewComment handles post requests to add a comment to a post
// It expects a multipart form with the comment content and an optional image
func NewComment(w http.ResponseWriter, r *http.Request, userID int) {

	var newComment models.Comment

	err := ParseContent(r, &newComment)
	if err != nil {
		log.Println("Error parsing comment data:", err)
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Bad Request"})
		return
	}

	// If multipart form, get the form values
	if newComment.CommentContent == "" {
		newComment.CommentContent = r.FormValue("comment_content")
		newComment.PostID, err = strconv.Atoi(r.FormValue("post_id"))
		if err != nil {
			log.Println("Invalid postID: ", err)
			ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Invalid post ID"})
			return
		}
		newComment.CommentImage = SaveUploadedFile(r, "comment_image", "comment")

		postIDStr := r.FormValue("post_id")
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			log.Println("Invalid post_id:", err)
			ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Invalid post_id"})
			return
		}
		newComment.PostID = postID
	}

	if !database.ValidatePostID(newComment.PostID) {
		log.Println("Invalid postID: ", newComment.PostID)
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Invalid post ID"})
		return
	}

	if newComment.CommentContent != "" {
		// Insert comment into the database
		err := database.AddCommentIntoDB(newComment.PostID, userID, newComment.CommentContent, newComment.CommentImage)
		if err != nil {
			ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
			return
		}
	} else {
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Comment content cannot be empty"})
		return
	}

	ResponseHandler(w, http.StatusOK, models.Response{Message: "Comment added to database"})
}
