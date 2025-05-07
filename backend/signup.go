package backend

import (
	"encoding/json"
	"log"
	"net/http"
	"net/mail"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// signUp handles both GET and POST requests for user registration
func SignUp(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleSignUpPost(w, r)
	default:
		ResponseHandler(w, http.StatusMethodNotAllowed, "Method Not Allowed")
	}
}

func handleSignUpPost(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON body into the LoginData struct
	var signUpData SignUpData
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
	if !IsValidUsername(signUpData.Username) {
		status = http.StatusBadRequest
		message = "Invalid username: must be 3-20 characters, letters, numbers, or _"
	} else if !isValidEmail(signUpData.Email) {
		status = http.StatusBadRequest
		message = "Invalid email address"
	} else if signUpData.Password == "" {
		status = http.StatusBadRequest
		message = "Password cannot be empty"
	} else if signUpData.Age == "" {
		status = http.StatusBadRequest
		message = "Please enter your age"
	} else if signUpData.Gender == "" {
		status = http.StatusBadRequest
		message = "Gender is still missing"
	} else if signUpData.LastName == "" || signUpData.FirstName == "" {
		status = http.StatusBadRequest
		message = "Please enter your first and last name"
	} else {
		uniqueUsername, uniqueEmail, err := isUsernameOrEmailUnique(signUpData.Username, signUpData.Email)
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
		err = insertUserIntoDB(
			signUpData.Username,
			signUpData.Age,
			signUpData.Gender,
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

// isUsernameOrEmailUnique checks if the username or email is unique in the database
func isUsernameOrEmailUnique(username, email string) (bool, bool, error) {
	username = strings.ToLower(username)
	email = strings.ToLower(email)

	var count int
	err := db.QueryRow(`
        SELECT COUNT(*) 
        FROM User 
        WHERE username = ?`, username).Scan(&count)
	if err != nil || count != 0 {
		return false, false, err
	}
	err = db.QueryRow(`
        SELECT COUNT(*) 
        FROM User 
        WHERE email = ?`, email).Scan(&count)
	if err != nil || count != 0 {
		return true, false, err
	}
	return true, true, nil // Returns true if neither username nor email exists
}
