package app

import (
	"encoding/json"
	"net/http"
	"social_network/database"
	"social_network/models"
	"strconv"
)

// ProfileResponse contains all data needed for a profile page
type ProfileResponse struct {
	User           models.User   `json:"user"`
	IsOwnProfile   bool          `json:"is_own_profile"`
	Posts          []models.Post `json:"posts"`
	FollowersCount int           `json:"followers_count"`
	FollowingCount int           `json:"following_count"`
}

// When user clicks any profile
func ServeProfile(w http.ResponseWriter, r *http.Request) {
	profileIDStr := r.URL.Query().Get("user_id")
	if profileIDStr == "" {
		http.Error(w, "(Error1 in ServeProfile)User ID is required", http.StatusBadRequest)
		return
	}
	profileID, err := strconv.Atoi(profileIDStr)
	if err != nil {
		http.Error(w, "(Error2 in ServeProfile)Invalid User ID", http.StatusBadRequest)
		return
	}

	isLoggedIn, viewerID := VerifySession(r)

	profileUser, err := database.GetUser(profileID)
	if err != nil {
		http.Error(w, "(Error3 in ServeProfile)User not found", http.StatusNotFound)
		return
	}

	// Check if viewer is the profile owner
	isOwnProfile := isLoggedIn && viewerID == profileID

	// Check privacy - if it's not public and not the owner, deny access
	if !isOwnProfile && !profileUser.IsPublic {
		http.Error(w, "(Error4 in ServeProfile)This profile is private", http.StatusForbidden)
		return
	}

	// Get posts, including privacy filter
	posts, err := database.GetUserPosts(profileID, viewerID, isOwnProfile)
	if err != nil {
		http.Error(w, "(Error5 in ServeProfile)Error fetching posts", http.StatusInternalServerError)
		return
	}

	// Get followers and following counts
	followersCount, _ := database.GetFollowersCount(profileID)
	followingCount, _ := database.GetFollowingCount(profileID)

	// Create the response object
	response := ProfileResponse{
		User:           profileUser,
		IsOwnProfile:   isOwnProfile,
		Posts:          posts,
		FollowersCount: followersCount,
		FollowingCount: followingCount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func ChangeProfileVisibility() {

	//Toggle between public and private in Users table
}

func DisplayFollowData() {

	// Check if following or followers
	// Get all ids from Followers table according to above specification
	// Get all names and pictures from the Users table and return
}
