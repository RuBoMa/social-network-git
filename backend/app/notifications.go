package app

import (
	"log"
	"net/http"
	"social_network/database"
)

// ServeUnreadNotifications retrieves unread notifications for a user
// It expects a userID to be passed in the request context
func ServeUnreadNotifications(w http.ResponseWriter, r *http.Request, userID int) {
	notifications, err := database.GetUnreadNotifications(userID)
	if err != nil {
		log.Println("Error fetching notifications:", err)
		ResponseHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	ResponseHandler(w, http.StatusOK, notifications)
}

// MarkNotificationRead marks a notification as read
// It expects a JSON body with the notification ID and read status
// If the notification ID is valid and the status is true, it marks the notification as read
// func MarkNotificationRead(w http.ResponseWriter, r *http.Request) {
// 	var notification models.Notification

// 	err := ParseContent(r, &notification)
// 	if err != nil {
// 		log.Println("Error parsing notification:", err)
// 		ResponseHandler(w, http.StatusBadRequest, "Bad Request")
// 		return
// 	}

// 	isValid, err := database.IsValidNotificationID(notification.NotificationID)
// 	if err != nil {
// 		log.Println("Error validating notification ID:", err)
// 		ResponseHandler(w, http.StatusInternalServerError, "Internal Server Error")
// 		return
// 	} else if !isValid {
// 		log.Println("Invalid notification ID:", notification.NotificationID)
// 		ResponseHandler(w, http.StatusBadRequest, "Invalid Notification ID")
// 		return
// 	} else {
// 		if notification.IsRead {
// 			err := database.NotificationSeen(notification.NotificationID)
// 			if err != nil {
// 				log.Println("Error marking notification as read:", err)
// 				ResponseHandler(w, http.StatusInternalServerError, "Internal Server Error")
// 				return
// 			}
// 			ResponseHandler(w, http.StatusOK, "Notification marked as read")
// 		} else {
// 			ResponseHandler(w, http.StatusAlreadyReported, "Notification status unchanged")
// 		}
// 	}
// }

/*
Using this style we should be able to call the function like this:
app.SendNotification("follow_request", receiverID, &requestID, nil)
app.SendNotification("group_invite", receiverID, &requestID, nil)
app.SendNotification("event_created", memberID, nil, &eventID)
If we want to add more types of notifications, we can just add them to the switch statement in the database function
*/
