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

	var response models.ProfileResponse
	var err error

	if userID < 1 {
		log.Println("Error: user_id not provided")
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Bad Request"})
		return
	}

	isLoggedIn, viewerID := VerifySession(r)

	response.User, err = database.GetUser(userID)
	if err != nil {
		log.Println("Error fetching user profile:", err)
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Bad Request"})
		return
	}

	// Check if viewer is the profile owner
	response.IsOwnProfile = isLoggedIn && viewerID == userID

	response.IsFollower = false

	if isLoggedIn && !response.IsOwnProfile {
		followers, err := database.GetFollowers(userID)
		if err != nil {
			log.Println("Error fetching followers:", err)
		} else {
			for _, follower := range followers {
				if follower.UserID == viewerID {
					response.IsFollower = true
					break
				}
			}
		}
	}

	if response.IsOwnProfile && !response.User.IsPublic {
		response.FollowRequests, err = database.GetOwnFollowRequests(userID)
		if err != nil {
			log.Println("Error fetching follow requests:", err)
			ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
			return
		}
	}

	// Get posts, based on the profile's privacy settings
	if response.IsOwnProfile || response.User.IsPublic || response.IsFollower {
		response.Posts, err = database.GetUserPosts(userID, viewerID, response.IsOwnProfile)
		if err != nil {
			log.Println("Error fetching posts:", err)
			ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
			return
		}
	}

	if !response.IsOwnProfile {
		// Check if the user is followed by or following the viewer
		canMessage, err := database.IsEitherOneFollowing(viewerID, userID)
		if err != nil {
			log.Println("Error checking following status:", err)
			ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
			return
		}
		// Check if there is a existing conversation between the viewer and this user
		hasConversation, err := database.HasExistingConversation(viewerID, userID)
		if err != nil {
			log.Println("Error checking existing conversation:", err)
			ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
			return
		}

		response.ShowChatButton = canMessage && !hasConversation

	}

	// Check if user has requested to follow this profile
	response.HasRequested = false
	if isLoggedIn && !response.IsOwnProfile {
		response.HasRequested, err = database.HasPendingFollowRequest(viewerID, userID)
		if err != nil {
			log.Println("Error checking follow request:", err)
		}
	}
	// Check requests table, is there this userID as sender and profileUser.UserID as receiver and status is "pending"
	// If found, set hasRequested to true

	// Get followers and following counts
	response.FollowersCount, _ = database.GetFollowersCount(userID)
	response.FollowingCount, _ = database.GetFollowingCount(userID)

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
	responseData := map[string]interface{}{
		"message":   "Profile visibility updated successfully",
		"is_public": newVisibility,
	}
	ResponseHandler(w, http.StatusOK, responseData)
}

func DisplayFollowData() {

	// Check if following or followers
	// Get all ids from Followers table according to above specification
	// Get all names and pictures from the Users table and return
}
