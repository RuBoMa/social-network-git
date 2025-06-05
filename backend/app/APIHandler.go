package app

import (
	"log"
	"net/http"
	"social_network/app/chat"
	"social_network/database"
	"social_network/models"
	"strconv"
	"strings"
)

// APIHandler is the main handler for API requests
// It parses the request, verifies the session, and routes the request to the appropriate handler
func APIHandler(w http.ResponseWriter, r *http.Request) {

	route := ParseRoute(r)
	if route.Err != nil {
		log.Println("Error parsing route:", route.Err)
		ResponseHandler(w, http.StatusNotFound, "Invalid URL")
		return
	}

	log.Println("Handling request for page:", route.Page)
	// Verify the session and check if the user is logged in
	loggedIn, userID := VerifySession(r)

	if !loggedIn && route.Page != "login" && route.Page != "signup" {
		log.Println("Unauthorized access attempt to:", route.Page)
		ResponseHandler(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Handle different routes based on the URL path
	switch r.Method {

	case http.MethodGet:

		switch route.Page {
		case "feed":
			HandleFeed(w, r, userID, route.GroupID)
		case "post":
			HandlePostGet(w, r, route.PostID, userID)
		case "profile":
			ServeProfile(w, r, route.ProfileID)
		case "my-groups":
			ServeUsersGroups(w, r, userID)
		case "all-groups":
			ServeAllGroups(w, r)
		case "group":
			if route.SubAction == "" {
				ServeGroup(w, r, route.GroupID, userID)
			} else if route.SubAction == "invite" {
				ServeNonGroupMembers(w, r, route.GroupID)
			} else {
				ResponseHandler(w, http.StatusNotFound, "Page Not Found")
				return
			}
		case "event":
			ServeEvent(w, r, route.EventID, userID)
		case "followers", "following":
			var id int
			if route.ProfileID != 0 {
				id = route.ProfileID
			} else {
				id = userID
			}
			if route.Page == "followers" {
				GetFollowers(w, id)
			} else {
				GetFollowing(w, id)
			}
		case "users":
			ServeUsers(w, r)
		case "notifications":
			ServeUnreadNotifications(w, r, userID)
		case "search":
			log.Println("Search query:", route.SearchParam)
			Search(w, r, route.SearchParam, userID)
		default:
			ResponseHandler(w, http.StatusNotFound, "Page Not Found")
			return
		}

	case http.MethodPost:

		switch route.Page {
		case "comment":
			NewComment(w, r, userID)
		case "login":
			HandleLogin(w, r)
		case "signup":
			HandleSignUp(w, r)
		case "create-post":
			NewPost(w, r, userID)
		case "create-group":
			CreateGroup(w, r, userID)
		case "create-event":
			CreateGroupEvent(w, r, userID)
		case "logout":
			Logout(w, r, userID)
			chat.CloseConnection(userID)
		case "request":
			log.Println("Request received.")
			HandleRequests(w, r, userID)
		case "event":
			CreateGroupEvent(w, r, userID)
		case "event-attendance":
			MarkEventAttendance(w, r, userID)
		case "privacy":
			UpdateProfilePrivacy(w, r, userID)
		default:
			ResponseHandler(w, http.StatusNotFound, "Page Not Found")
			return
		}

	default:
		ResponseHandler(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}
}

// ParseRoute parses the URL path and query parameters to extract route information
// It returns a RouteInfo struct containing the page, post ID, and any errors encountered
func ParseRoute(r *http.Request) models.RouteInfo {
	// Filter URL path
	parts := strings.Split(r.URL.Path, "/")
	var filtered []string
	for _, p := range parts {
		if p != "" && p != "api" {
			filtered = append(filtered, p)
		}
	}

	if len(filtered) == 0 {
		return models.RouteInfo{Err: http.ErrNoLocation}
	}

	info := models.RouteInfo{Page: filtered[0]}
	query := r.URL.Query()

	if len(filtered) > 1 {
		info.SubAction = filtered[1]
	}

	// Validate the query parameters
	if qParam := query.Get("q"); qParam != "" {
		info.SearchParam = qParam
	} else if eventIDStr := query.Get("event_id"); eventIDStr != "" {
		if id, err := strconv.Atoi(eventIDStr); err == nil {
			info.EventID = id
		} else {
			info.Err = err
			return info
		}
	} else if postIDStr := query.Get("post_id"); postIDStr != "" {
		if id, err := strconv.Atoi(postIDStr); err == nil {
			valid := database.ValidatePostID(id)
			if !valid {
				log.Println("Invalid postID: ", id)
				info.Err = http.ErrNoLocation
				return info
			}
			info.PostID = id
		} else {
			info.Err = err
			return info
		}
	} else if userIDStr := query.Get("user_id"); userIDStr != "" {
		if id, err := strconv.Atoi(userIDStr); err == nil {
			info.ProfileID = id
		} else {
			info.Err = err
			return info
		}
	} else if groupIDStr := query.Get("group_id"); groupIDStr != "" {
		if id, err := strconv.Atoi(groupIDStr); err == nil {
			info.GroupID = id
		} else {
			info.Err = err
			return info
		}
	}

	// Validate that the page has needed parameters
	if (info.Page == "group" && info.GroupID == 0) ||
		(info.Page == "post" && info.PostID == 0) ||
		(info.Page == "profile" && info.ProfileID == 0) ||
		(info.Page == "search" && info.SearchParam == "") ||
		(info.Page == "event" && info.EventID == 0) {
		info.Err = http.ErrNoLocation
		info.Page = ""
	}

	return info
}
