package app

import (
	"encoding/json"
	"log"
	"net/http"
	"net/mail"
	"regexp"
	"social_network/database"
	"social_network/models"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// HandleSignUp handles the user sign-up form submission
// It validates the input, hashes the password, and stores the user in the database
func HandleSignUp(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON body into the LoginData struct
	err := r.ParseMultipartForm(10 << 20) // max 10MB
	if err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	data := models.User{
		Nickname:    r.FormValue("nickname"),
		DateOfBirth: r.FormValue("date_of_birth"), // need validation
		FirstName:   r.FormValue("first_name"),
		LastName:    r.FormValue("last_name"),
		Email:       r.FormValue("email"),
		Password:    r.FormValue("password"),
		AboutMe:     r.FormValue("about_me"), // need validation

		IsPublic: r.FormValue("is_public") == "on",
	}

	data.AvatarPath = SaveUploadedFile(r, "avatar", "profile")

	status := http.StatusCreated
	message := "Login successful"

	// Validate username
	if !IsValidUsername(data.Nickname) {
		status = http.StatusBadRequest
		message = "Invalid username: must be 3-20 characters, letters, numbers, or _"
	} else if !IsValidEmail(data.Email) {
		status = http.StatusBadRequest
		message = "Invalid email address"
	} else if data.Password == "" {
		status = http.StatusBadRequest
		message = "Password cannot be empty"
	} else if data.DateOfBirth == "" {
		status = http.StatusBadRequest
		message = "Please enter your date of birth"
	} else if data.LastName == "" || data.FirstName == "" {
		status = http.StatusBadRequest
		message = "Please enter your first and last name"
	} else {
		uniqueEmail, err := database.IsEmailUnique(data.Email)
		if err != nil {
			log.Println("Error checking if email is unique:", err)
			ResponseHandler(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		if !uniqueEmail {
			status = http.StatusConflict
			message = "Email is already registered to existing user"
		}
	}

	if message == "Login successful" && status == http.StatusCreated {
		// Hash the password
		hashedPassword, err := HashPassword(data.Password)
		if err != nil {
			log.Println("Error hashing password:", hashedPassword)
			ResponseHandler(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		// Insert user into database
		err = database.AddUserIntoDB(
			data.Email,
			hashedPassword,
			data.FirstName,
			data.LastName,
			data.DateOfBirth,
			data.AvatarPath,
			data.Nickname,
			data.AboutMe,
			data.IsPublic,
		)
		if err != nil {
			log.Println("Error inserting user into database:", err)
			ResponseHandler(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	ResponseHandler(w, status, message)
}

// HandleLoginPost handles the user login form submission
func HandleLogin(w http.ResponseWriter, r *http.Request) {

	// Decode the JSON body into the LoginData struct
	var loginData models.LoginData
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loginData)
	if err != nil {
		log.Println("Error decoding the login data")
		ResponseHandler(w, http.StatusBadRequest, "Bad Request")
		return
	}

	message := "Login successful"
	status := http.StatusOK
	userID, hashedPassword, err := database.GetUserCredentials(loginData.Email)
	if err != nil {
		log.Println("Invalid username")
		status = http.StatusUnauthorized
		message = "Invalid username or email"
	} else {
		err := VerifyPassword(hashedPassword, loginData.Password)
		if err != nil {
			log.Println("Invalid password")
			status = http.StatusUnauthorized
			message = "Invalid password"
		} else {
			// Create session
			if err := CreateSession(w, r, userID); err != nil {
				log.Println("Error creating session")
				ResponseHandler(w, http.StatusInternalServerError, "Internal Server Error")
				return
			}
		}
	}

	if status == http.StatusOK {
		username, err := database.GetUsername(userID)
		if err != nil {
			log.Println("Error getting username")
			ResponseHandler(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		response := models.User{
			Nickname: username,
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			ResponseHandler(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		return
	}

	ResponseHandler(w, status, message)
}

// Logout logs out the user by deleting the session from the database and setting the session cookie to expire
func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Println("Session cookie not found:", err)
		ResponseHandler(w, http.StatusBadRequest, "No session cookie found")
		return
	}

	err = database.DeleteActiveSession(cookie.Value)
	if err != nil {
		// if return error, it didn't find any sessions to delete
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
	})

	ResponseHandler(w, http.StatusOK, "Logout successful")
}

// Authenticate checks if the user is logged in by verifying the session ID
func Authenticate(w http.ResponseWriter, loggedIn bool, userID int) {
	status := http.StatusUnauthorized
	message := "No current sessions"

	if loggedIn {
		status = http.StatusOK
		message = strconv.Itoa(userID)
	}

	ResponseHandler(w, status, message)
}

// hashPassword hashes the user's password using bcrypt
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

// verifyPassword compares the hashed password with the password provided by the user
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// isValidEmail checks if the email address is valid
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	regex := regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+$`)
	return regex.MatchString(email)
}

// IsValidUsername checks if the username is valid
func IsValidUsername(username string) bool {
	if len(username) == 0 {
		return true
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`) // Only letters, numbers, and _
	return re.MatchString(username)
}
