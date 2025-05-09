package app

import (
	"encoding/json"
	"net/http"
	"social_network/models"
	"strconv"

	"social_network/database"
)

// When user clicks any profile
func ServeProfile(w http.ResponseWriter, r *http.Request) {
	profileIDStr := r.URL.Query().Get("user_id")
	if profileIDStr == "" {
		http.Error(w, "(Error1 in ServerProfile)User ID is required", http.StatusBadRequest)
		return
	}
	profileID, err := strconv.Atoi(profileIDStr)
	if err != nil {
		http.Error(w, "(Error2 in ServerProfile)Invalid User ID", http.StatusBadRequest)
		return
	}

	isLoggedIn, viewerID := VerifySession(r)

	profileUser, err := database.GetUsername(profileID)
	if err != nil {
		http.Error(w, "(Error3 in ServerProfile)User not found", http.StatusNotFound)
		return
	}

	isOwner := isLoggedIn && viewerID == profileID

	if !isOwner && !profileUser.IsPublic {
		http.Error(w, "(Error4 in ServerProfile)Profile is private", http.StatusForbidden)
		return
	}

	posts, err := FetchFeed(profileID, isOwner)
	if err != nil {
		http.Error(w, "(Error5 in ServerProfile)Error fetching posts", http.StatusInternalServerError)
		return
	}
	followersCount, err := database.GetFollowersCount(profileID)
	followingCount, err := database.GetFollowingCount(profileID)

	type ProfileResponse struct {
		UserID         int      `json:"user_id"`
		isOwner        bool     `json:"is_own_profile"`
		Posts          []string `json:"posts"`
		FollowersCount int      `json:"followers_count"`
		FollowingCount int      `json:"following_count"`
	}

	var postTitles []string
	for _, post := range posts {
		postTitles = append(postTitles, post.PostTitle)
	}
	response := ProfileResponse{
		User: models.User{
			UserID:     profileUser.UserID,
			Nickname:   profileUser.Nickname,
			FirstName:  profileUser.FirstName,
			LastName:   profileUser.LastName,
			Email:      profileUser.Email,
			AboutMe:    profileUser.AboutMe,
			AvatarPath: profileUser.AvatarPath,
			IsPublic:   profileUser.IsPublic,
		},
		isOwner:        isOwner,
		Posts:          postTitles,
		FollowersCount: followersCount,
		FollowingCount: followingCount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	//Check if it's their own profile (compare user_ids)
	// If yes, add that information to the response

	//Check if the profile is public or private and add to responce

	// Get profile information (minus password)
	// Get all posts by the user (inlude privacy) Only titles?
	// Get following SUM (user_id == follower_id) and followers SUM (user_id == followed_id)

}

func ChangeProfileVisibility() {

	//Toggle between public and private in Users table
}

func DisplayFollowData() {

	// Check if following or followers
	// Get all ids from Followers table according to above specification
	// Get all names and pictures from the Users table and return
}
