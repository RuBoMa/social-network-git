package app

import (
	"log"
	"net/http"
	"social_network/database"
)

func ServeUnreadNotifications(w http.ResponseWriter, r *http.Request, userID int) {
	notifications, err := database.GetUnreadNotifications(userID)
	if err != nil {
		log.Println("Error fetching notifications:", err)
		ResponseHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	log.Println("Notifications", notifications)
	ResponseHandler(w, http.StatusOK, notifications)
}

/*
Using this style we should be able to call the function like this:
app.SendNotification("follow_request", receiverID, &requestID, nil)
app.SendNotification("group_invite", receiverID, &requestID, nil)
app.SendNotification("event_created", memberID, nil, &eventID)
If we want to add more types of notifications, we can just add them to the switch statement in the database function
*/
