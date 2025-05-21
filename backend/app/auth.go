package app

import (
	"log"
	"net/http"
	"net/mail"
	"regexp"
	"social_network/database"
	"social_network/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// HandleSignUp handles the user sign-up form submission
// It validates the input, hashes the password, and stores the user in the database
func HandleSignUp(w http.ResponseWriter, r *http.Request) {

	data := models.User{}

	err := ParseContent(r, &data)
	if err != nil {
		log.Println("Error decoding the sign-up data")
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Invalid form"})
		return
	}

	// If multipart form, get the form values
	if data.Email == "" {
		data.Nickname = r.FormValue("nickname")
		data.DateOfBirth = r.FormValue("date_of_birth")
		data.FirstName = r.FormValue("first_name")
		data.LastName = r.FormValue("last_name")
		data.Email = r.FormValue("email")
		data.Password = r.FormValue("password")
		data.AboutMe = r.FormValue("about_me")
		data.IsPublic = r.FormValue("is_public") == "true"

		data.AvatarPath = SaveUploadedFile(r, "avatar", "profile")
	}

	log.Println("Parsed sign-up data:", data)

	status := http.StatusCreated
	message := models.Response{
		Message: "Login successful",
	}

	// Validate username
	if !IsValidUsername(data.Nickname) {
		status = http.StatusBadRequest
		message.Message = "Invalid username: must be 3-20 characters, letters, numbers, or _"
	} else if !IsValidEmail(data.Email) {
		status = http.StatusBadRequest
		message.Message = "Invalid email address"
	} else if data.Password == "" {
		status = http.StatusBadRequest
		message.Message = "Password cannot be empty"
	} else if !IsValidDateOfBirth(data.DateOfBirth) {
		status = http.StatusBadRequest
		message.Message = "Please enter your date of birth"
	} else if data.LastName == "" || data.FirstName == "" {
		status = http.StatusBadRequest
		message.Message = "Please enter your first and last name"
	} else if !IsValidAboutMe(data.AboutMe) {
		status = http.StatusBadRequest
		message.Message = "About me must be between 10 and 500 characters and cannot contain disallowed HTML tags"
	} else {
		uniqueEmail, err := database.IsEmailUnique(data.Email)
		if err != nil {
			log.Println("Error checking if email is unique:", err)
			ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
			return
		}
		if !uniqueEmail {
			status = http.StatusConflict
			message.Message = "Email is already registered to existing user"
		}
	}

	if message.Message == "Login successful" && status == http.StatusCreated {
		// Hash the password
		hashedPassword, err := HashPassword(data.Password)
		if err != nil {
			log.Println("Error hashing password:", err)
			ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
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
			ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
			return
		}
	}

	ResponseHandler(w, status, message)
}

// HandleLoginPost handles the user login form submission
func HandleLogin(w http.ResponseWriter, r *http.Request) {

	// Decode the JSON body into the LoginData struct
	var loginData models.LoginData

	err := ParseContent(r, &loginData)
	if err != nil {
		log.Println("Error decoding the login data")
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "Bad Request"})
		return
	}

	message := models.Response{Message: "Login successful"}
	status := http.StatusOK
	userID, hashedPassword, err := database.GetUserCredentials(loginData.Email)
	if err != nil {
		log.Println("Invalid email: ", err)
		status = http.StatusUnauthorized
		message.Message = "Invalid email"
	} else {
		err := VerifyPassword(hashedPassword, loginData.Password)
		if err != nil {
			log.Println("Invalid password")
			status = http.StatusUnauthorized
			message.Message = "Invalid password"
		} else {
			// Create session
			if err := CreateSession(w, r, userID); err != nil {
				log.Println("Error creating session:", err)
				ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
				return
			}
		}
	}

	if status == http.StatusOK {
		var user models.User
		user.UserID = userID
		var err error
		user, err = database.GetUser(userID)
		if err != nil {
			log.Println("Error getting user info:", err)
			ResponseHandler(w, http.StatusInternalServerError, models.Response{Message: "Internal Server Error"})
			return
		}

		ResponseHandler(w, status, user)
		return
	}

	ResponseHandler(w, status, message)
}

// Logout logs out the user by deleting the session from the database and setting the session cookie to expire
func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Println("Session cookie not found:", err)
		ResponseHandler(w, http.StatusBadRequest, models.Response{Message: "No session cookie found"})
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

	ResponseHandler(w, http.StatusOK, models.Response{Message: "Logout successful"})
}

// Authenticate checks if the user is logged in by verifying the session ID
func Authenticate(w http.ResponseWriter, loggedIn bool, userID int) {

	if loggedIn {
		ResponseHandler(w, http.StatusOK, models.User{UserID: userID})
	} else {
		ResponseHandler(w, http.StatusUnauthorized, models.Response{Message: "No current sessions"})
	}

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

func IsValidDateOfBirth(dob string) bool {
	// Check if the date is in the format DD-MM-YYYY
	_, err := time.Parse("2006-01-02", dob)
	return err == nil
}

func IsValidAboutMe(about string) bool {
	if len(about) == 0 {
		return true
	}

	if len(about) > 500 {
		return false
	}
	// Check for disallowed HTML tags
	disallowed := regexp.MustCompile(`(?i)<(script|iframe|embed|object|style)>`)
	return !disallowed.MatchString(about)
}
