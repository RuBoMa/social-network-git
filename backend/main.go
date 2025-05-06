package main

import (
	"real-time-forum/database"
	"real-time-forum/server"
)

func main() {

	// Initialize database
	db := database.InitDB()
	defer db.Close()

	database.MakeTables()

	server.Run()
}
