package app

import (
	"social_network/database"
	"social_network/models"
)

func FetchFeed(userID int) ([]models.Post, error) {
	posts, err := database.GetPosts(userID)
	if err != nil {
		return nil, err
	}

	return posts, nil
}
