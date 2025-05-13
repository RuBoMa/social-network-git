package app

import (
	"log"
	"net/http"
	"social_network/database"
	"social_network/models"
)

// ServeProfile handles requests to view a user's profile
// It retrieves the user's information, posts, and followers/following counts
func ServeProfile(w http.ResponseWriter, r *http.Request) {
	var profile models.User
	err := ParseContent(r, &profile)
	if err != nil {
		log.Println("(Error1 in ServeProfile)Error parsing the profile data:", err)
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Bad Request"})
		return
	}

	if profile.UserID < 1 {
		log.Println("(Error2 in ServeProfile)Error: user_id not provided")
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Bad Request"})
		return
	}

	isLoggedIn, viewerID := VerifySession(r)

	profileUser, err := database.GetUser(profile.UserID)
	if err != nil {
		log.Println("(Error3 in ServeProfile)Error fetching user profile:", err)
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Bad Request"})
		return
	}

	// Check if viewer is the profile owner
	isOwnProfile := isLoggedIn && viewerID == profile.UserID

	// Check privacy - if it's not public and not the owner, deny access
	// NEEDS WORK, we don't deny access. We have to check if the viewer is follower and has access
	// If not, we still need to return basic profile information
	// Added IsFollower (bool) to the ProfileResponse struct
	isFollower := false

	if isLoggedIn && !isOwnProfile {
		followers, err := database.GetFollowing(profile.UserID)
		if err != nil {
			log.Println("(Error4 in ServeProfile)Error fetching followers:", err)
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
		posts, err = database.GetUserPosts(profile.UserID, viewerID, isOwnProfile)
		if err != nil {
			log.Println("(Error5 in ServeProfile)Error fetching posts:", err)
			ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
			return
		}
	} else {
		// If the profile is private and the viewer is not a follower - return empty posts
		posts = []models.Post{}
	}

	// Get followers and following counts
	followersCount, _ := database.GetFollowersCount(profile.UserID)
	followingCount, _ := database.GetFollowingCount(profile.UserID)

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
