package app

import (
	"log"
	"net/http"
	"social_network/database"
	"social_network/models"
	"strconv"
)

// FetchFeed retrieves the posts for a given user
func FetchFeed(userID int) ([]models.Post, error) {
	posts, err := database.GetPosts(userID)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

// HandlePostGet handles get requests to a specific post
// It retrieves the post details and comments associated with it
func HandlePostGet(w http.ResponseWriter, r *http.Request, postID, userID int) {
	post, err := database.GetPostDetails(postID)
	if err != nil {
		log.Println("Error fetching post details:", err)
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
		return
	}

	ResponseHandler(w, http.StatusOK, post)

}

// NewPost handles post requests to create a new post
// It expects a multipart form with the post title, content, privacy setting, and an optional image
// It also associates the post with a group if provided
func NewPost(w http.ResponseWriter, r *http.Request, userID int) {

	var newPost models.Post

	err := r.ParseMultipartForm(10 << 20) // limit: 10MB
	if err != nil {
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Invalid multipart form"})
		return
	}

	newPost.PostTitle = r.FormValue("post_title")
	newPost.PostContent = r.FormValue("post_content")
	newPost.Privacy = r.FormValue("privacy")
	groupIDSstr := r.FormValue("group_id")

	if groupIDSstr != "" {
		groupID, err := strconv.Atoi(groupIDSstr)
		if err != nil {
			ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Invalid group ID"})
			return
		}
		newPost.Group.GroupID = groupID
	}

	if newPost.PostTitle == "" || newPost.PostContent == "" || newPost.Privacy == "" {
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Title, content or privacy setting cannot be empty"})
		return
	}

	newPost.PostImage = SaveUploadedFile(r, "post_image", "posts")

	err = database.AddPostIntoDB(newPost.PostTitle, newPost.PostContent, newPost.PostImage, newPost.Privacy, userID, newPost.Group.GroupID)
	if err != nil {
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
		return
	}

	ResponseHandler(w, http.StatusOK, models.Response{Message: "Message added to database"})

}

// NewComment handles post requests to add a comment to a post
// It expects a multipart form with the comment content and an optional image
func NewComment(w http.ResponseWriter, r *http.Request, postID, userID int) {

	var newComment models.Comment

	err := r.ParseMultipartForm(10 << 20) // limit: 10MB
	if err != nil {
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Invalid multipart form"})
		return
	}

	newComment.CommentContent = r.FormValue("comment_content")
	newComment.CommentImage = SaveUploadedFile(r, "comment_image", "comments")

	if newComment.CommentContent != "" {
		// Insert comment into the database
		err := database.AddCommentIntoDB(postID, userID, newComment.CommentContent, newComment.CommentImage)
		if err != nil {
			ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
			return
		}
	}

	HandlePostGet(w, r, postID, userID)
}
