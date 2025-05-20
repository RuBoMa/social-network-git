package app

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// ParseContent function parses the content of the HTTP request
// It checks the Content-Type (multipart/form-data or application/json) header and decodes the request body accordingly
func ParseContent(r *http.Request, content interface{}) error {
	ct := r.Header.Get("Content-Type")

	if strings.HasPrefix(ct, "multipart/form-data") {
		err := r.ParseMultipartForm(10 << 20) // max 10MB
		if err != nil {
			log.Println("Error parsing multipart form:", err)
			return err
		}
	} else if strings.HasPrefix(ct, "application/json") {
		err := json.NewDecoder(r.Body).Decode(content)
		if err != nil {
			log.Println("Error decoding JSON:", err)
			return err
		}
	}

	return nil
}
