package app

import (
	"log"
	"net/http"
	"social_network/database"
	"social_network/models"
)

// ServeProfile handles requests to view a user's profile
// It retrieves the user's information, posts, and followers/following counts
func ServeProfile(w http.ResponseWriter, r *http.Request, userID int) {

	if userID < 1 {
		log.Println("Error: user_id not provided")
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Bad Request"})
		return
	}

	isLoggedIn, viewerID := VerifySession(r)

	profileUser, err := database.GetUser(userID)
	if err != nil {
		log.Println("Error fetching user profile:", err)
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Bad Request"})
		return
	}

	// Check if viewer is the profile owner
	isOwnProfile := isLoggedIn && viewerID == userID

	isFollower := false

	if isLoggedIn && !isOwnProfile {
		followers, err := database.GetFollowing(userID)
		if err != nil {
			log.Println("Error fetching followers:", err)
			isFollower = false
		} else {
			for _, followerID := range followers {
				if followerID == viewerID {
					isFollower = true
					break
				}
			}
		}
	}

	// Get posts, based on the profile's privacy settings
	var posts []models.Post
	if isOwnProfile || profileUser.IsPublic || isFollower {
		posts, err = database.GetUserPosts(userID, viewerID, isOwnProfile)
		if err != nil {
			log.Println("Error fetching posts:", err)
			ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
			return
		}
	} else {
		// If the profile is private and the viewer is not a follower - return empty posts
		posts = []models.Post{}
	}

	// Get followers and following counts
	followersCount, _ := database.GetFollowersCount(userID)
	followingCount, _ := database.GetFollowingCount(userID)

	response := models.ProfileResponse{
		User:           profileUser,
		IsOwnProfile:   isOwnProfile,
		IsFollower:     isFollower,
		Posts:          posts,
		FollowersCount: followersCount,
		FollowingCount: followingCount,
	}

	ResponseHandler(w, http.StatusOK, response)

}

func ChangeProfileVisibility() {

	//Toggle between public and private in Users table
}

func DisplayFollowData() {

	// Check if following or followers
	// Get all ids from Followers table according to above specification
	// Get all names and pictures from the Users table and return
}
