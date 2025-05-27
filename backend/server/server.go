package server

import (
	"log"
	"net/http"
	"social_network/app"
	"social_network/app/chat"
	"strings"
)

func Run() {
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))
	// One API Handler for api calls
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle CORS preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		ct := r.Header.Get("Content-Type")

		if strings.HasPrefix(ct, "application/json") || strings.HasPrefix(ct, "multipart/form-data") || ct == "" {
			APIHandler(w, r)
			return
		}
		log.Println("Unsupported Content-Type:", ct)
		app.ResponseHandler(w, http.StatusUnsupportedMediaType, "Unsupported Content-Type")
	})

	// Handler for chat
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		log.Println("WebSocket connection attempt")

		HandleConnections(w, r)
	})

	// Start message broadcaster
	go chat.BroadcastMessages()

	log.Println("Message broadcaster started")

	log.Println("Server is running on http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting the server:", err)
	}
}
