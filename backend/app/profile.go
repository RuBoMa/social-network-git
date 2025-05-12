package app

import (
	"log"
	"net/http"
	"social_network/database"
	"social_network/models"
	"strconv"
)

// ServeProfile handles requests to view a user's profile
// It retrieves the user's information, posts, and followers/following counts
func ServeProfile(w http.ResponseWriter, r *http.Request) {
	profileIDStr := r.URL.Query().Get("user_id")
	if profileIDStr == "" {
		log.Println("Error: user_id not provided")
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Bad Request"})
		return
	}
	profileID, err := strconv.Atoi(profileIDStr)
	if err != nil {
		log.Println("Error converting user_id to int:", err)
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Bad Request"})
		return
	}

	isLoggedIn, viewerID := VerifySession(r)

	profileUser, err := database.GetUser(profileID)
	if err != nil {
		log.Println("Error fetching user profile:", err)
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Bad Request"})
		return
	}

	// Check if viewer is the profile owner
	isOwnProfile := isLoggedIn && viewerID == profileID

	// Check privacy - if it's not public and not the owner, deny access
	// NEEDS WORK, we don't deny access. We have to check if the viewer is follower and has access
	// If not, we still need to return basic profile information
	// Added IsFollower (bool) to the ProfileResponse struct
	if !isOwnProfile && !profileUser.IsPublic {
		http.Error(w, "(Error4 in ServeProfile)This profile is private", http.StatusForbidden)
		return
	}

	// Get posts, including privacy filter
	posts, err := database.GetUserPosts(profileID, viewerID, isOwnProfile)
	if err != nil {
		log.Println("Error fetching posts:", err)
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
		return
	}

	// Get followers and following counts
	followersCount, _ := database.GetFollowersCount(profileID)
	followingCount, _ := database.GetFollowingCount(profileID)

	response := models.ProfileResponse{
		User:           profileUser,
		IsOwnProfile:   isOwnProfile,
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
