package main

import (
	"real-time-forum/database"
	"real-time-forum/server"
)

func main() {

	// Initialize database
	database.InitDB()
	defer func() {
		database.CloseDB()
	}()

	server.Run()
}
