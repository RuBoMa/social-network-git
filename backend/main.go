package main

import (
	"social_network/database"
	"social_network/server"
)

func main() {

	// Initialize database
	database.InitDB()
	defer func() {
		database.CloseDB()
	}()

	server.Run()
}
