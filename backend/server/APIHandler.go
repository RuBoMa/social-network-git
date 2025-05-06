package server

import (
	"log"
	"net/http"
	"real-time-forum/app"
	"real-time-forum/database"
	"strconv"
	"strings"
)

func APIHandler(w http.ResponseWriter, r *http.Request) {

	route := ParseRoute(r.URL.Path)
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

	// Handle different routes based on the URL path
	loggedIn, userID := app.VerifySession(r)

	switch r.Method {

	case http.MethodGet:

		switch route.Page {
		case "feed":
			HandleFeed(w, r) // Returns posts to be shown in feed
		case "auth":
			app.Authenticate(w, loggedIn, userID)
		case "post":
			app.HandlePostPageGet(w, r, route.PostID, userID)
		case "create-post":
			app.FetchCategories(w, r)
		default:
			app.ResponseHandler(w, http.StatusNotFound, "Page Not Found")
			return
		}

	case http.MethodPost:

		switch route.Page {
		case "post":
			if route.SubAction == "vote" {
				app.HandleVote(w, r, route.PostID, userID)
			} else if route.SubAction == "comment" {
				app.HandleComment(w, r, route.PostID, userID)
			} else {
				app.ResponseHandler(w, http.StatusBadRequest, "Bad Request")
			}
		case "login":
			app.HandleLoginPost(w, r)
		case "signup":
			app.HandleSignUpPost(w, r)
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

func ParseRoute(path string) RouteInfo {
	parts := strings.Split(path, "/")
	var filtered []string
	for _, p := range parts {
		if p != "" && p != "api" {
			filtered = append(filtered, p)
		}
	}

	info := RouteInfo{}
	if len(filtered) > 0 {
		info.Page = filtered[0]
	}

	if info.Page == "post" {
		if len(filtered) > 1 {
			id, err := strconv.Atoi(filtered[1])
			if err != nil {
				info.Err = err
				return info
			}
			info.PostID = id

			if len(filtered) >= 3 {
				info.SubAction = filtered[2]
			}
		} else {
			info.Page = ""
		}
	}
	return info
}
