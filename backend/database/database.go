package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3" // SQLite3 driver; the blank import ensures the driver is included "_" is important!!!
)

var db *sql.DB

// InitDB initializes the SQLite database and returns a database connection object
func InitDB() {
	var err error
	db, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Println("Failed to initialize the database")
		log.Fatal(err)

	}
}

func CloseDB() {
	db.Close()

	ApplyMigrations(db)
}

// applies the database migrations
func ApplyMigrations(db *sql.DB) {
	// create migration driver instance using the existing sql.DB connection for SQLite
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		log.Fatal("Error creating migration driver:", err)
	}
	// create new migrate instane using driver and path to migration files
	mig, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		"sqlite3",
		driver,
	)
	if err != nil {
		log.Fatal("Error creating migrate instance:", err)
	}

	// Apply migrations
	err = mig.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("Error applying migrations:", err)
	} else if err == migrate.ErrNoChange {
		fmt.Println("no new migrations to apply")
	} else {
		fmt.Println("Migrations applied successfully")
	}
}
