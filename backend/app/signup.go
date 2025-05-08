package app

import (
	"encoding/json"
	"log"
	"net/http"
	"net/mail"
	"regexp"
	"social_network/database"
	"social_network/models"

	"golang.org/x/crypto/bcrypt"
)

func HandleSignUpPost(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON body into the LoginData struct
	var signUpData models.SignUpData
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&signUpData)
	if err != nil {
		log.Println("error decoding the data")
		ResponseHandler(w, http.StatusBadRequest, "Bad Request")
		return
	}

	status := http.StatusCreated
	message := "Login successful"

	// Validate username
	if !IsValidUsername(signUpData.Nickname) {
		status = http.StatusBadRequest
		message = "Invalid username: must be 3-20 characters, letters, numbers, or _"
	} else if !isValidEmail(signUpData.Email) {
		status = http.StatusBadRequest
		message = "Invalid email address"
	} else if signUpData.Password == "" {
		status = http.StatusBadRequest
		message = "Password cannot be empty"
	} else if signUpData.DateOfBirth == "" {
		status = http.StatusBadRequest
		message = "Please enter your age"
	} else if signUpData.LastName == "" || signUpData.FirstName == "" {
		status = http.StatusBadRequest
		message = "Please enter your first and last name"
	} else {
		uniqueUsername, uniqueEmail, err := database.IsUsernameOrEmailUnique(signUpData.Nickname, signUpData.Email)
		if err != nil {
			log.Println("Error checking if username is unique:", err)
			ResponseHandler(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		if !uniqueUsername {
			status = http.StatusConflict
			message = "Username is already taken"
		} else if !uniqueEmail {
			status = http.StatusConflict
			message = "Email is already registered to existing user"
		}
	}

	if message == "Login successful" && status == http.StatusCreated {
		// Hash the password
		hashedPassword, err := hashPassword(signUpData.Password)
		if err != nil {
			log.Println("Error hashing password:", err)
			ResponseHandler(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		// Insert user into database
		err = database.InsertUserIntoDB(
			signUpData.Nickname,
			signUpData.DateOfBirth,
			signUpData.FirstName,
			signUpData.LastName,
			signUpData.Email,
			hashedPassword,
		)
		if err != nil {
			log.Println("Error inserting user into database:", err)
			ResponseHandler(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	ResponseHandler(w, status, message)
}

// hashPassword hashes the user's password using bcrypt
func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

// isValidEmail checks if the email address is valid
func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	regex := regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+$`)
	return regex.MatchString(email)
}

// IsValidUsername checks if the username is valid
func IsValidUsername(username string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`) // Only letters, numbers, and _
	return re.MatchString(username)
}
