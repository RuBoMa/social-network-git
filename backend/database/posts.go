package database

import (
	"fmt"
	"log"
	"social_network/models"
	"strconv"
)

// GetPosts fetches all posts from the database and returns them as a slice of PostDetails
func GetPostIDs(userID int) ([]models.PostDetails, error) {
	var posts []models.PostDetails

	// Query to get all posts ordered by creation date
	query := `
		SELECT Post.id
		FROM Post
		ORDER BY Post.created_at DESC;
	`
	rows, err := db.Query(query)
	if err != nil {
		log.Println("Error fetching posts:", err)
		return nil, err
	}
	defer rows.Close()

	// Loop through the rows and fetch details for each post
	for rows.Next() {
		var postID int
		if err := rows.Scan(&postID); err != nil {
			log.Println("Error scanning post ID:", err)
			return nil, err
		}

		// Get the details for each post
		post, err := GetPostDetails(postID, userID)
		if err != nil {
			log.Println("Error getting post details:", err)
			return nil, err
		}

		posts = append(posts, *post)
	}

	// Return the list of posts
	return posts, nil
}

// // HandleCategory converts the category ID into a string and returns validated ID
func HandleCategory(category string) (int, error) {

	categoryID, err := strconv.Atoi(category)
	if err != nil {
		log.Println("Error converting categoryID", err)
		return 0, err
	}

	valid := ValidateCategoryID(categoryID)
	if !valid {
		log.Println("Invalid categoryID", category)
		return 0, fmt.Errorf("invalid category id: %s", category)
	}

	return categoryID, nil

}

// // ValidateCategoryID checks if the category ID given exists in the databse
func ValidateCategoryID(categoryID int) bool {
	var category int
	err := db.QueryRow("SELECT id FROM Category WHERE id = ?", categoryID).Scan(&category)
	if err != nil {
		log.Println("Error scanning category ID:", err)
		return false
	}
	return true
}
