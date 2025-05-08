package server

import (
	"log"
	"net/http"
	"social_network/app"
)

func Run() {

	// One API Handler for api calls
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			app.ResponseHandler(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
			return
		}
		APIHandler(w, r)
	})

	// Handler for chat
	http.HandleFunc("/ws", HandleConnections)

	// Start message broadcaster
	go BroadcastMessages()

	log.Println("Server is running on http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting the server:", err)
	}
}
