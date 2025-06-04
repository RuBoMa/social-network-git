package seed

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	FirstName   string
	LastName    string
	Nickname    string
	Email       string
	Password    string
	DateOfBirth string
	AvatarPath  string
	AboutMe     string
	IsPublic    bool
}

func SeedUsers(db *sql.DB) error {
	// Define the SQL statement to insert a user
	hashed, err := bcrypt.GenerateFromPassword([]byte("12345"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	dob := time.Date(1995, 1, 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02 15:04:05")

	users := []User{
		{
			FirstName:   "Fanni",
			LastName:    "Dev",
			Nickname:    "fanni",
			Email:       "fanni@test.com",
			Password:    string(hashed),
			DateOfBirth: dob,
			AvatarPath:  "",
			AboutMe:     "Just a dev",
			IsPublic:    true,
		},
		{
			FirstName:   "Alex",
			LastName:    "Coder",
			Nickname:    "alex",
			Email:       "alex@test.com",
			Password:    string(hashed),
			DateOfBirth: dob,
			AvatarPath:  "",
			AboutMe:     "Love coding",
			IsPublic:    false,
		},
		{
			FirstName:   "Roope",
			LastName:    "Hongisto",
			Nickname:    "roopelainen",
			Email:       "roope@test.com",
			Password:    string(hashed),
			DateOfBirth: dob,
			AvatarPath:  "",
			AboutMe:     "Help me, I'm a dev",
			IsPublic:    false,
		},
	}

	sqlStatement := `
	INSERT INTO Users (first_name, last_name, nickname, email, password_hash, date_of_birth, avatar_path, about_me, is_public)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT DO NOTHING;`

	for _, user := range users {
		_, err = db.Exec(sqlStatement,
			user.FirstName, user.LastName, user.Nickname, user.Email, user.Password,
			user.DateOfBirth, user.AvatarPath, user.AboutMe, user.IsPublic,
		)
		if err != nil {
			return err
		}
	}

	return nil

}
