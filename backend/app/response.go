package app

import (
	"encoding/json"
	"net/http"
	"social_network/models"
)

// ResponseHandler function handles the HTTP response
// It sets the appropriate headers and encodes the response message as JSON
// Cross-Origin Resource Sharing (CORS) headers are set to allow requests from frontend
func ResponseHandler(w http.ResponseWriter, statusCode int, data interface{}) {

	w.Header().Set("Content-Type", "application/json")

	if msg, ok := data.(string); ok {
		data = models.Response{Message: msg}
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
