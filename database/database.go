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

// InitDB initializes the SQLite database and returns a database connection object
func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Println("Failed to initialize the database")
		log.Fatal(err)

	}

	ApplyMigrations(db)
	return db
}

// applies the database migrations
func ApplyMigrations(db *sql.DB) {

	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		log.Fatal("Error creating migration driver:", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		"sqlite3",
		driver,
	)
	if err != nil {
		log.Fatal("Error creating migrate instance:", err)
	}

	// Apply migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("Error applying migrations:", err)
	} else if err == migrate.ErrNoChange {
		fmt.Println("no new migrations to apply")
	} else {
		fmt.Println("Migrations applied successfully")
	}
}
