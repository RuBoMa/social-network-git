package server

import (
	"log"
	"net/http"
	"social_network/app"
	"strings"
)

func Run() {

	// One API Handler for api calls
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")

		if strings.HasPrefix(ct, "application/json") || strings.HasPrefix(ct, "multipart/form-data") {
			APIHandler(w, r)
			return
		}
		app.ResponseHandler(w, http.StatusUnsupportedMediaType, "Unsupported Content-Type")
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
