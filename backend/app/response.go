package app

import (
	"encoding/json"
	"net/http"
	"social_network/models"
)

func ResponseHandler(w http.ResponseWriter, statusCode int, message string) {
	response := models.Response{Message: message}

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
