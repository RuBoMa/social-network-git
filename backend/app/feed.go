package app

import (
	"real-time-forum/backend/database"
	"real-time-forum/backend/models"
)

func FetchFeed(userID int) ([]models.PostDetails, error) {
	posts, err := database.GetPostIDs(userID)
	if err != nil {
		return nil, err
	}

	return posts, nil
}
