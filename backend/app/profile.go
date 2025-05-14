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

func ChangeProfileVisibility(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		ResponseHandler(w, http.StatusMethodNotAllowed, models.Response{
			Message: "Method not allowed",
		})
		return
	}

	// Verify the user is logged in
	isLoggedIn, userID := VerifySession(r)
	if !isLoggedIn || userID < 1 {
		ResponseHandler(w, http.StatusUnauthorized, models.Response{
			Message: "Unauthorized",
		})
		return
	}

	// Get current user profile to check current visibility
	user, err := database.GetUser(userID)
	if err != nil {
		log.Println("Error fetching user profile:", err)
		ResponseHandler(w, http.StatusInternalServerError, models.Response{
			Message: "Failed to retrieve user profile",
		})
		return
	}

	// Toggle visibility
	newVisibility := !user.IsPublic

	// Update the visibility in the database
	err = database.UpdateProfileVisibility(userID, newVisibility)
	if err != nil {
		log.Println("Error updating profile visibility:", err)
		ResponseHandler(w, http.StatusInternalServerError, models.Response{
			Message: "Failed to update profile visibility",
		})
		return
	}

	// Return success response with new visibility status
	ResponseHandler(w, http.StatusOK, models.Response{
		Message: "Profile visibility updated successfully",
		Data: map[string]bool{
			"is_public": newVisibility,
		},
	})
}

func DisplayFollowData() {

	// Check if following or followers
	// Get all ids from Followers table according to above specification
	// Get all names and pictures from the Users table and return
}
