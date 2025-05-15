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
		app.ResponseHandler(w, http.StatusNotFound, "Invalid URL")
		return
	}

	if route.PostID > 0 {
		valid := database.ValidatePostID(route.PostID)
		if !valid {
			log.Println("Invalid postID: ", route.PostID)
			app.ResponseHandler(w, http.StatusNotFound, "Page Not Found")
			return
		}
	}

	loggedIn, userID := app.VerifySession(r)

	if route.Page == "profile" {
		if route.ProfileID == 0 {
			route.ProfileID = userID
		} else if !database.IsValidUserID(route.ProfileID) {
			log.Println("Invalid profileID: ", route.ProfileID)
			app.ResponseHandler(w, http.StatusNotFound, "Page Not Found")
			return

		}
	}

	// Handle different routes based on the URL path

	switch r.Method {

	case http.MethodGet:

		switch route.Page {
		case "feed":
			app.HandleFeed(w, r, userID) // Returns posts to be shown in feed
		case "auth":
			app.Authenticate(w, loggedIn, userID)
		case "post":
			app.HandlePostGet(w, r, route.PostID, userID)
		case "profile":
			app.ServeProfile(w, r, route.ProfileID)
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
		case "logout":
			app.Logout(w, r)
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
	path := r.URL.Path
	parts := strings.Split(path, "/")
	var filtered []string
	for _, p := range parts {
		if p != "" && p != "api" {
			filtered = append(filtered, p)
		}
	}

	info := models.RouteInfo{}
	if len(filtered) > 0 {
		info.Page = filtered[0]
	}

	if info.Page == "post" {
		postIDStr := r.URL.Query().Get("post_id")
		if postIDStr == "" {
			info.Page = ""
			return info
		}

		id, err := strconv.Atoi(postIDStr)
		if err != nil {
			info.Err = err
			return info
		}
		info.PostID = id
	} else if info.Page == "profile" {
		userIDStr := r.URL.Query().Get("user_id")
		if userIDStr != "" {
			id, err := strconv.Atoi(userIDStr)
			if err != nil {
				info.Err = err
				return info
			}
			info.ProfileID = id
		}
	}

	return info
}
