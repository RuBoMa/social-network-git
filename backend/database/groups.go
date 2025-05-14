package database

import (
	"social_network/models"
	"strings"
	"time"
)

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

// AddGroupIntoDB adds a new group to the database
// It takes a models.Group object as input and inserts it into the Groups table
func AddGroupIntoDB(group models.Group) (int, error) {
	result, err := db.Exec("INSERT INTO Groups_Table (title, description, creator_id, created_at) VALUES (?, ?, ?, ?)",
		group.GroupName, group.GroupDesc, group.GroupCreator.UserID, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return 0, err
	}

	groupID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(groupID), nil
}

// IsGroupNameUnique checks if the given title is unique in the database
func IsGroupNameUnique(title string) (bool, error) {
	title = strings.ToLower(title)

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Groups_Table WHERE title = ?", title).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

// IsValidGroupID checks if the given group ID exists in the database
func IsValidGroupID(groupID int) bool {
	if groupID == 0 {
		return false
	}
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Groups_Table WHERE id = ?", groupID).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}

// AddEventIntoDB adds a new event to the database
// It takes a models.Event object as input and inserts it into the Events table
func AddEventIntoDB(event models.Event) (int, error) {
	result, err := db.Exec("INSERT INTO Events (group_id, creator_id, title, description, event_time, created_at) VALUES (?, ?, ?, ?, ?)",
		event.Group.GroupID, event.CreatorID, event.Title, event.Description, event.EventDate, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return 0, err
	}

	eventID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(eventID), nil
}

func IsValidEventID(eventID int) bool {
	if eventID == 0 {
		return false
	}
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Events WHERE id = ?", eventID).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}

func AddEventResponseIntoDB(response models.EventResponse) (int, error) {
	result, err := db.Exec("INSERT INTO Events_Responses (event_id, user_id, response, created_at) VALUES (?, ?, ?, ?)",
		response.EventID, response.UserID, response.Response, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return 0, err
	}

	responseID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(responseID), nil
}
