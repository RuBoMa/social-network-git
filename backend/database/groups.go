package database

import "social_network/models"

func GetGroupMembers(groupID int) ([]models.User, error) {
	var users []models.User

	rows, err := db.Query(`
		SELECT u.id
		FROM Users u
		JOIN Group_Members gm ON u.id = gm.user_id
		WHERE gm.group_id = ?`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserID); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
