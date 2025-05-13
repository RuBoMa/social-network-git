package app

import (
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

// GroupInvitation handles group invitations
// It parses the request body to get sender, receiver, and group details,
// and checks if the sender and receiver are valid users
func GroupInvitations(w http.ResponseWriter, r *http.Request) {
	invitation := models.Request{}

	err := ParseContent(r, &invitation)
	if err != nil {
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Invalid form"})
		return
	}

	// Validate group ID, sender ID, and receiver ID from the database
	// Validation checks also if the id is missing
	if !database.IsValidGroupID(invitation.Group.GroupID) || !database.IsValidUserID(invitation.Sender.UserID) || !database.IsValidUserID(invitation.Receiver.UserID) {
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Invalid group ID, sender ID or receiver ID"})
		return
	}
	// Add group invitation to the database as pending
	invitation.RequestID, err = database.AddGroupInvitationIntoDB(invitation)
	if err != nil {
		ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Database error"})
		return
	}
	// Send group invitation to the frontend
	ResponseHandler(w, http.StatusOK, invitation)

}

func AnswerToGroupInvitation() {

}
