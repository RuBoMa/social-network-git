package backend

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Login handles both GET and POST requests for user authentication
func Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		HandleLoginPost(w, r)
	default:
		ResponseHandler(w, http.StatusMethodNotAllowed, "Method Not Allowed")
	}
}

// HandleLoginPost handles the user login form submission
func HandleLoginPost(w http.ResponseWriter, r *http.Request) {

	// Decode the JSON body into the LoginData struct
	var loginData LoginData
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loginData)
	if err != nil {
		log.Println("Error decoding the login data")
		ResponseHandler(w, http.StatusBadRequest, "Bad Request")
		return
	}

	message := "Login successful"
	status := http.StatusOK
	userID, hashedPassword, err := getUserCredentials(loginData.Username)
	if err != nil {
		log.Println("Invalid username")
		status = http.StatusUnauthorized
		message = "Invalid username or email"
	} else {
		err := verifyPassword(hashedPassword, loginData.Password)
		if err != nil {
			log.Println("Invalid password")
			status = http.StatusUnauthorized
			message = "Invalid password"
		} else {
			// Create session
			if err := CreateSession(w, userID); err != nil {
				log.Println("Error creating session")
				ResponseHandler(w, http.StatusInternalServerError, "Internal Server Error")
				return
			}
		}
	}

	if status == http.StatusOK {
		username, err := GetUsername(userID)
		if err != nil {
			log.Println("Error getting username")
			ResponseHandler(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		response := SignUpData{
			Username: username,
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

// verifyPassword compares the hashed password with the password provided by the user
func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
