package app

import (
	"log"
	"net/http"
	"social_network/database"
	"social_network/models"
	"strings"
)

// ServeGroups handles the request to get all groups for the groupBar
func ServeAllGroups(w http.ResponseWriter, r *http.Request) {
	var groups []models.Group
	var err error
	groups, err = database.GetAllGroups()
	if err != nil {
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Database error"})
		return
	}

	ResponseHandler(w, http.StatusOK, groups)
}

// ServeUsersGroups handles the request to get all groups the user is a member of
func ServeUsersGroups(w http.ResponseWriter, r *http.Request, userID int) {
	var groups []models.Group
	var err error
	groups, err = database.GetUsersGroups(userID)
	if err != nil {
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Database error"})
		return
	}

	ResponseHandler(w, http.StatusOK, groups)
}

// ServeGroup handles the request to get a specific group for the group page
func ServeGroup(w http.ResponseWriter, r *http.Request, groupID, userID int) {
	var group models.Group
	var err error

	// Check if the group ID is valid
	if !database.IsValidGroupID(groupID) {
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Invalid group ID"})
		return
	}

	group, err = database.GetGroupByID(groupID)
	if err != nil {
		log.Println("Error retrieving group by ID:", err)
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Database error"})
		return
	}

	group.GroupCreator, err = database.GetUser(group.GroupCreator.UserID)
	if err != nil {
		log.Println("Error retrieving group creator:", err)
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Database error"})
		return
	}

	group.GroupMembers, err = database.GetGroupMembers(groupID)
	if err != nil {
		log.Println("Error retrieving group members:", err)
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Database error"})
	}

	group.IsMember = false
	for _, member := range group.GroupMembers {
		if member.UserID == userID {
			group.IsMember = true
			break
		}
	}

	group.RequestStatus, group.RequestID, err = database.ActiveRequest(userID, groupID)
	if err != nil {
		log.Println("Error retrieving request status:", err)
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Database error"})
		return
	}

	group.GroupEvents, err = database.GetGroupEvents(groupID)
	if err != nil {
		log.Println("Error retrieving group events:", err)
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Database error"})
		return
	}

	ResponseHandler(w, http.StatusOK, group)
}

// ServeGroupRequests handles the request to get all requests for a specific group
func ServeGroupRequests(w http.ResponseWriter, r *http.Request, groupID int) {
	var requests []models.Request
	var err error

	// Check if the group ID is valid
	if !database.IsValidGroupID(groupID) {
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Invalid group ID"})
		return
	}

	requests, err = database.GetGroupRequests(groupID)
	if err != nil {
		log.Println("Error retrieving group requests:", err)
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Database error"})
		return
	}

	ResponseHandler(w, http.StatusOK, requests)
}

// CreateGroup handles the creation of a new group
// It parses the request body to get group details, checks for uniqueness of group name,
// and adds the group to the database
func CreateGroup(w http.ResponseWriter, r *http.Request, userID int) {
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

	group.GroupCreator, err = database.GetUser(userID)
	if err != nil {
		log.Println("Error retrieving group creator:", err)
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Database error"})
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

	// Add the creator as the first member of the group
	err = database.AddGroupMemberIntoDB(group.GroupID, group.GroupCreator.UserID)
	if err != nil {
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Database error"})
		return
	}

	// Return group information so that frontend can show it
	ResponseHandler(w, http.StatusOK, group.GroupID)
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
	err = database.AddNotificationIntoDB(models.NotifGroupInvite, request, models.Event{})
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

	if request.Status == "accepted" {
		// Add the user to the group if the request is accepted
		err = database.AddGroupMemberIntoDB(request.Group.GroupID, request.Sender.UserID)
		if err != nil {
			ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Database error"})
			return
		}
	}

	// HOW TO HANDLE THE NOTIFICATION?

	ResponseHandler(w, http.StatusOK, request)

}

// GROUP EVENTS
// Create Group Event --> notification
// Going/Not Going

// CreateGroupEvent handles the creation of a new group event
// It parses the request body to get event details, and adds the event to the database
// It also adds an unread notification to the group members
func CreateGroupEvent(w http.ResponseWriter, r *http.Request, userID int) {
	event := models.Event{}
	err := ParseContent(r, &event)
	if err != nil {
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Invalid form"})
		return
	}

	event.Title = strings.TrimSpace(event.Title)
	event.Description = strings.TrimSpace(event.Description)
	log.Println("Event:", event)
	if event.Title == "" || event.Description == "" || event.EventDate == "" {
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Event title, description and date is required"})
		return
	}
	// VALIDATE EVENT DATE !?!?

	if !database.IsValidGroupID(event.Group.GroupID) {
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Invalid group ID"})
		return
	}

	event.Creator.UserID = userID

	event.EventID, err = database.AddEventIntoDB(event)
	if err != nil {
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Database error"})
		return
	}

	err = database.AddNotificationIntoDB(models.NotifEventCreated, models.Request{}, event)
	if err != nil {
		log.Println("Error saving notification:", err)
		// Currently not crashing the server if notification fails
	}

	ResponseHandler(w, http.StatusOK, event)
}

// MarkEventAttendance handles the response to a group event
// It parses the request body to get event ID and response (going/not going)
func MarkEventAttendance(w http.ResponseWriter, r *http.Request, userID int) {
	answer := models.EventResponse{}
	err := ParseContent(r, &answer)
	if err != nil {
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Invalid form"})
		return
	}
	if answer.Event.EventID == 0 || answer.Response == "" {
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Event ID and response is required"})
		return
	}

	if answer.Response != "going" && answer.Response != "not going" {
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Invalid response"})
		return
	}

	if !database.IsValidEventID(answer.Event.EventID) {
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Invalid event ID"})
		return
	}

	answer.User.UserID = userID

	answer.ResponseID, err = database.AddEventResponseIntoDB(answer)
	if err != nil {
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Database error"})
		return
	}

	ResponseHandler(w, http.StatusOK, answer)

}

func ServeEvent(w http.ResponseWriter, r *http.Request, eventID, userID int) {
	var event models.Event
	var err error

	// Check if the event ID is valid
	if !database.IsValidEventID(eventID) {
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Invalid event ID"})
		return
	}

	event, err = database.GetEventByID(eventID)
	if err != nil {
		log.Println("Error retrieving event by ID:", err)
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Database error"})
		return
	}

	event.Creator, err = database.GetUser(event.Creator.UserID)
	if err != nil {
		log.Println("Error retrieving event creator:", err)
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Database error"})
		return
	}

	ResponseHandler(w, http.StatusOK, event)
}
