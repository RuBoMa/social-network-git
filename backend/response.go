package backend

import (
	"encoding/json"
	"net/http"
)

func ResponseHandler(w http.ResponseWriter, statusCode int, message string) {
	response := Response{Message: message}

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
