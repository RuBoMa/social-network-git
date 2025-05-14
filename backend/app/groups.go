package app

import (
	"log"
	"net/http"
	"social_network/database"
	"social_network/models"
	"strings"
)

// CreateGroup handles the creation of a new group
// It parses the request body to get group details, checks for uniqueness of group name,
// and adds the group to the database
func CreateGroup(w http.ResponseWriter, r *http.Request) {
	group := models.Group{}
	err := ParseContent(r, &group)
	if err != nil {
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Invalid form"})
		return
	}

	group.GroupName = strings.TrimSpace(group.GroupName)
	group.GroupDesc = strings.TrimSpace(group.GroupDesc)

	if group.GroupName == "" || group.GroupDesc == "" {
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Group name and description is required"})
		return
	}

	// Check if group title already exists in the Groups table
	exists, err := database.IsGroupNameUnique(group.GroupName)
	if err != nil {
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Database error"})
		return
	}
	if !exists {
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Group name already exists"})
		return
	}
	// If not, add information to database
	group.GroupID, err = database.AddGroupIntoDB(group)
	if err != nil {
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Database error"})
		return
	}

	// Return group information so that frontend can show it
	ResponseHandler(w, http.StatusOK, group)
}

// JoinGroup handles group join requests
func JoinGroup(w http.ResponseWriter, r *http.Request, request models.Request) {

	//Check if the user is already a member of the group??

	if request.Status == "invited" || request.Status == "requested" {
		GroupRequests(w, r, request)
		return
	} else if request.Status == "accepted" || request.Status == "rejected" {
		// Handle group invitation response
		AnswerToGroupRequest(w, r, request)
	} else {
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Invalid status"})
	}

}

// GroupRequests handles group invitations/requests
// It saves the request and notification into the database
func GroupRequests(w http.ResponseWriter, r *http.Request, request models.Request) {
	var err error

	if request.Status == "invited" {
		if request.Receiver.UserID == 0 {
			ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "ReceiverID is missing from group invitation"})
			return
		}
	}

	// Add group invitation to the database with current status
	request.RequestID, err = database.AddRequestIntoDB(request)
	if err != nil {
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Database error"})
		return
	}
	// Save notification into database
	err = database.AddNotificationIntoDB(request, models.Event{})
	if err != nil {
		log.Println("Error saving notification:", err)
		// Currently not crashing the server if notification fails
	}

	ResponseHandler(w, http.StatusOK, request)

}

// AnswerToGroupRequest handles the response to a group invitation/request
// It updates the status of the group invitation in the database
func AnswerToGroupRequest(w http.ResponseWriter, r *http.Request, request models.Request) {
	// Check if the request ID is valid
	if request.RequestID == 0 {
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Invalid request ID"})
		return
	}

	// Update the status of the group invitation in the database
	err := database.UpdateRequestStatus(request)
	if err != nil {
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Database error"})
		return
	}

	// HOW TO HANDLE THE NOTIFICATION?

	ResponseHandler(w, http.StatusOK, request)

}

// GROUP EVENTS
// Create Group Event
// Going/Not Going
