package app

import (
	"encoding/json"
	"net/http"
)

// ResponseHandler function handles the HTTP response
// It sets the appropriate headers and encodes the response message as JSON
// Cross-Origin Resource Sharing (CORS) headers are set to allow requests from frontend
func ResponseHandler(w http.ResponseWriter, statusCode int, data interface{}) {

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // check with front their port
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true") // needed for cookies ??
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
