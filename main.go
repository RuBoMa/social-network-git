package main

import (
	"real-time-forum/backend/database"
	"real-time-forum/backend/server"
)

func main() {

	// Initialize database
	db := database.InitDB()
	defer db.Close()

	database.MakeTables()

	server.Run()
}
