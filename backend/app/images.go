package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// SaveUploadedFile saves an uploaded file to a specified directory and returns the relative path
// for database storage. It creates the directory if it doesn't exist and generates a unique filename.
// The file is saved in the format: /uploads/{subDir}/{timestamp}.{extension}
func SaveUploadedFile(r *http.Request, formKey, subDir string) string {
	file, handler, err := r.FormFile(formKey)
	if err != nil {
		return "" // nothing uploaded
	}
	defer file.Close()

	// Ensure directory exists
	saveDir := filepath.Join("../uploads", subDir)
	err = os.MkdirAll(saveDir, os.ModePerm)
	if err != nil {
		return ""
	}

	// Generate unique filename
	ext := filepath.Ext(handler.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	fullPath := filepath.Join(saveDir, filename)

	// Create and copy file
	dst, err := os.Create(fullPath)
	if err != nil {
		return ""
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return ""
	}

	// Return relative path for database
	return "/uploads/" + subDir + "/" + filename
}
