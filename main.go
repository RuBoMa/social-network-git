package main

import (
	"log"
	"net/http"
	"real-time-forum/backend"
	"real-time-forum/database"
	"real-time-forum/websocket"
	"text/template"
)

func main() {

	// Parse and serve the template
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatal("Error parsing template")
		return
	}

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))

	// Initialize database
	db := database.InitDB()
	defer db.Close()

	database.MakeTables(db)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
		}
	})

	// One API Handler for api calls
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			backend.ResponseHandler(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
			return
		}
		backend.APIHandler(w, r, db)
	})

	// Handler for chat
	http.HandleFunc("/ws", websocket.HandleConnections)

	// Start message broadcaster
	go websocket.BroadcastMessages()

	log.Println("Server is running on http://localhost:8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting the server:", err)
	}
}
