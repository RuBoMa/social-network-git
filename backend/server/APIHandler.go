package server

import (
	"log"
	"net/http"
	"social_network/app"
	"social_network/database"
	"social_network/models"
	"strconv"
	"strings"
)

func APIHandler(w http.ResponseWriter, r *http.Request) {

	route := ParseRoute(r)
	if route.Err != nil {
		log.Println("Error parsing route:", route.Err)
		app.ResponseHandler(w, http.StatusNotFound, "Invalid URL")
		return
	}

	// log.Println("Parsed route:", route)

	loggedIn, userID := app.VerifySession(r)

	// Handle different routes based on the URL path

	switch r.Method {

	case http.MethodGet:

		switch route.Page {
		case "feed":
			app.HandleFeed(w, r, userID, route.GroupID) // Returns posts to be shown in feed
		case "auth":
			app.Authenticate(w, loggedIn, userID)
		case "post":
			app.HandlePostGet(w, r, route.PostID, userID)
		case "profile":
			app.ServeProfile(w, r, route.ProfileID)
		case "my-groups":
			app.ServeUsersGroups(w, r, userID)
		case "all-groups":
			app.ServeAllGroups(w, r)
		case "group":
			if route.SubAction == "" {
				app.ServeGroup(w, r, route.GroupID, userID)
			} else if route.SubAction == "invite" {
				app.ServeNonGroupMembers(w, r, route.GroupID)
			} else {
				app.ResponseHandler(w, http.StatusNotFound, "Page Not Found")
				return
			}
		case "event":
			app.ServeEvent(w, r, route.EventID, userID)
		case "followers", "following":
			var id int
			if route.ProfileID != 0 {
				id = route.ProfileID
			} else {
				id = userID
			}
			if route.Page == "followers" {
				app.GetFollowers(w, id)
			} else {
				app.GetFollowing(w, id)
			}
		case "users":
			app.ServeUsers(w, r)
		case "notifications":
			app.ServeUnreadNotifications(w, r, userID)
		case "search":
			log.Println("Search query:", route.SearchParam)
			app.Search(w, r, route.SearchParam, userID)
		default:
			app.ResponseHandler(w, http.StatusNotFound, "Page Not Found")
			return
		}

	case http.MethodPost:

		switch route.Page {
		case "comment":
			app.NewComment(w, r, userID)
		case "login":
			app.HandleLogin(w, r)
		case "signup":
			app.HandleSignUp(w, r)
		case "create-post":
			app.NewPost(w, r, userID)
		case "create-group":
			app.CreateGroup(w, r, userID)
		case "create-event":
			app.CreateGroupEvent(w, r, userID)
		case "logout":
			app.Logout(w, r)
		case "request":
			app.HandleRequests(w, r, userID)
		case "event":
			app.CreateGroupEvent(w, r, userID)
		case "event-attendance":
			app.MarkEventAttendance(w, r, userID)
		case "notifications":
			app.MarkNotificationRead(w, r)
		default:
			app.ResponseHandler(w, http.StatusNotFound, "Page Not Found")
			return
		}

	default:
		app.ResponseHandler(w, http.StatusMethodNotAllowed, "Method Not Allowed")
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

	if qParam := query.Get("q"); qParam != "" {
		info.SearchParam = qParam
	}

	if eventIDStr := query.Get("event_id"); eventIDStr != "" {
		if id, err := strconv.Atoi(eventIDStr); err == nil {
			info.EventID = id
		} else {
			info.Err = err
			return info
		}
	}

	// Try parsing all possible IDs independently
	if postIDStr := query.Get("post_id"); postIDStr != "" {
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
	}

	if userIDStr := query.Get("user_id"); userIDStr != "" {
		if id, err := strconv.Atoi(userIDStr); err == nil {
			info.ProfileID = id
		} else {
			info.Err = err
			return info
		}
	}

	if groupIDStr := query.Get("group_id"); groupIDStr != "" {
		if id, err := strconv.Atoi(groupIDStr); err == nil {
			info.GroupID = id
		} else {
			info.Err = err
			return info
		}
	}

	if info.Page == "group" && info.GroupID == 0 {
		info.Err = http.ErrNoLocation
		info.Page = ""
	}

	if info.Page == "search" && info.SearchParam == "" {
		info.Err = http.ErrNoLocation
		info.Page = ""
	}

	return info
}
