package app

import (
	"log"
	"net/http"
	"social_network/app/chat"
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

func ServeNotification(notificationID int) {
	notification, err := database.GetNotificationByID(notificationID)
	if err != nil {
		log.Println("Error fetching notification:", err)
		return
	}
	
	chat.BroadcastNotification(notification)
}
