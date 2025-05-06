package app

import (
	"real-time-forum/database"
	"real-time-forum/models"
)

func FetchFeed(userID int) ([]models.PostDetails, error) {
	posts, err := database.GetPostIDs(userID)
	if err != nil {
		return nil, err
	}

	return posts, nil
}
