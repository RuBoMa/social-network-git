package database

import (
	"database/sql"
	"log"
	"social_network/models"
	"strings"
	"time"
)

// GetAllGroups retrieves all groups from the database
func GetAllGroups() ([]models.Group, error) {
	var groups []models.Group

	rows, err := db.Query("SELECT id, title, description, creator_id, created_at FROM Groups_Table")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var group models.Group
		if err := rows.Scan(&group.GroupID, &group.GroupName, &group.GroupDesc, &group.GroupCreator.UserID, &group.GroupCreatedAt); err != nil {
			return nil, err
		}

		group.GroupCreator, err = GetUser(group.GroupCreator.UserID)
		if err != nil {
			log.Println("Error retrieving group creator:", err)
			return nil, err
		}
		groups = append(groups, group)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

func GetUsersGroups(userID int) ([]models.Group, error) {
	var groups []models.Group

	rows, err := db.Query(`
		SELECT g.id, g.title, g.description
		FROM Groups_Table g
		JOIN Group_Members gm ON g.id = gm.group_id
		WHERE gm.user_id = ?`, userID)
	if err != nil {
		log.Println("Error retrieving user's groups:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var group models.Group
		if err := rows.Scan(&group.GroupID, &group.GroupName, &group.GroupDesc); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

// GetGroupByID retrieves a group by its ID from the database
func GetGroupByID(groupID int) (models.Group, error) {
	var group models.Group

	err := db.QueryRow("SELECT id, creator_id, title, description, created_at FROM Groups_Table WHERE id = ?",
		groupID).Scan(&group.GroupID, &group.GroupCreator.UserID, &group.GroupName, &group.GroupDesc, &group.GroupCreatedAt)
	if err != nil {
		log.Println("Error retrieving group by ID:", err)
		return models.Group{}, err
	}

	return group, nil
}

func IsGroupMember(userID int, groupID int) (bool, error) {
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM Group_Members
		WHERE user_id = ? AND group_id = ?`, userID, groupID).Scan(&count)
	if err != nil {
		log.Println("Error checking group membership:", err)
		return false, err
	}
	log.Println("Group membership count:", count)

	return count > 0, nil
}

func GetGroupMembers(groupID int) ([]models.User, error) {
	var users []models.User

	rows, err := db.Query(`
		SELECT u.id
		FROM Users u
		JOIN Group_Members gm ON u.id = gm.user_id
		WHERE gm.group_id = ?`, groupID)
	if err != nil {
		log.Println("Error retrieving group members:", err)

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

func AddGroupMemberIntoDB(groupID, userID int) error {

	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM Group_Members WHERE group_id = ? AND user_id = ?
		)
	`, groupID, userID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		log.Println("User is already a member of the group")
		return nil
	}

	_, err = db.Exec(`
		INSERT INTO Group_Members (group_id, user_id, joined_at)
		VALUES (?, ?, ?)`,
		groupID, userID, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return err
	}

	return nil
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
	result, err := db.Exec("INSERT INTO Events (group_id, creator_id, title, description, event_time, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		event.Group.GroupID, event.Creator.UserID, event.Title, event.Description, event.EventDate, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Println("Error inserting event into database:", err)
		return 0, err
	}

	eventID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(eventID), nil
}

// GetGroupEvents retrieves all events for a specific group from the database
// It takes a groupID as input and returns a slice of models.Event
func GetGroupEvents(groupID int) ([]models.Event, error) {
	var events []models.Event

	rows, err := db.Query(`
		SELECT id, title, description, event_time
		FROM Events
		WHERE group_id = ?`, groupID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No events found for group ID:", groupID)
			return events, nil
		}
		log.Println("Error retrieving group events:", err)
		return events, err
	}
	defer rows.Close()

	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.EventID, &event.Title, &event.Description, &event.EventDate); err != nil {
			return events, err
		}
		event.Group.GroupID = groupID
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return events, err
	}

	return events, nil
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

// AddEventResponseIntoDB adds a new event response to the database
// It check first if the response already exists for the user and event and changes it if it does
func AddEventResponseIntoDB(response models.EventResponse) (int, error) {

	// Check if user has already responded to the event
	var existingStatus string
	var responseID int

	err := db.QueryRow(`
		SELECT id, response
		FROM Events_Responses
		WHERE event_id = ? AND user_id = ?
	`, response.Event.EventID, response.User.UserID).Scan(&responseID, &existingStatus)

	if err != nil {
		if err == sql.ErrNoRows {
			result, err := db.Exec("INSERT INTO Events_Responses (event_id, user_id, response, created_at) VALUES (?, ?, ?, ?)",
				response.Event.EventID, response.User.UserID, response.Response, time.Now().Format("2006-01-02 15:04:05"))
			if err != nil {
				return 0, err
			}

			response, err := result.LastInsertId()
			if err != nil {
				return 0, err
			}

			return int(response), nil

		} else {
			return 0, err
		}

	} else {

		// Update the existing response
		_, err = db.Exec(`
			UPDATE Events_Responses
			SET response = ?, updated_at = ?
			WHERE event_id = ? AND user_id = ?`,
			response.Response, time.Now().Format("2006-01-02 15:04:05"),
			response.Event.EventID, response.User.UserID)
		if err != nil {
			log.Println("Error updating existing event response:", err)
			return 0, err
		}

		return responseID, nil

	}
}

// SearchGroups searches for max 10 groups by title in the database
func SearchGroups(searchTerm string) ([]models.Group, error) {
	var groups []models.Group

	rows, err := db.Query(`
		SELECT id, title, description, creator_id, created_at
		FROM Groups_Table
		WHERE title LIKE ?
		ORDER BY
			CASE
				WHEN title = ? THEN 0
				ELSE 1
			END,
			title ASC
		LIMIT 10
		`,
		"%"+searchTerm+"%",
		searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var group models.Group
		if err := rows.Scan(&group.GroupID, &group.GroupName, &group.GroupDesc, &group.GroupCreator.UserID, &group.GroupCreatedAt); err != nil {
			return nil, err
		}

		group.GroupCreator, err = GetUser(group.GroupCreator.UserID)
		if err != nil {
			log.Println("Error retrieving group creator:", err)
			return nil, err
		}
		groups = append(groups, group)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

// SearchEvents searches for max 10 events by title or description in the database
func SearchEvents(searchTerm string, userID int) ([]models.Event, error) {
	var events []models.Event

	rows, err := db.Query(`
		SELECT e.id, e.group_id, e.title, e.description, e.event_time
		FROM Events e
		WHERE e.title LIKE ? OR e.description LIKE ?
		ORDER BY
			CASE
				WHEN e.title = ? THEN 0
				ELSE 1
			END,
			e.title ASC
		LIMIT 10
	`,
		"%"+searchTerm+"%",
		"%"+searchTerm+"%",
		searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.EventID, &event.Group.GroupID, &event.Title, &event.Description, &event.EventDate); err != nil {
			return nil, err
		}
		members, err := GetGroupMembers(event.Group.GroupID)
		if err != nil {
			log.Println("Error retrieving group members:", err)
			return nil, err
		}

		for _, member := range members {
			if member.UserID == userID {
				events = append(events, event)
				break
			}
		}

	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

// GetEventByID retrieves an event by its ID from the database
// It takes an eventID as input and returns a models.Event object including the creator details
func GetEventByID(eventID int) (models.Event, error) {
	var event models.Event

	err := db.QueryRow("SELECT id, group_id, creator_id, title, description, event_time FROM Events WHERE id = ?",
		eventID).Scan(&event.EventID, &event.Group.GroupID, &event.Creator.UserID, &event.Title, &event.Description, &event.EventDate)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No event found with ID:", eventID)
			return event, nil // No event found
		}
		log.Println("Error retrieving event by ID:", err)
		return event, err
	}

	return event, nil
}

// GetAttendingMembers retrieves all users who are attending a specific event
// It takes an eventID as input and returns a slice of models.User
func GetAttendingMembers(eventID int) ([]models.User, error) {
	var users []models.User

	rows, err := db.Query(`
		SELECT u.id
		FROM Users u
		JOIN Events_Responses er ON u.id = er.user_id
		WHERE er.event_id = ? AND er.response = 'going'`, eventID)
	if err != nil {
		log.Println("Error retrieving event attendees:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserID); err != nil {
			return nil, err
		}
		user, err = GetUser(user.UserID)
		if err != nil {
			log.Println("Error retrieving user:", err)
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// GetEventAttendance retrieves the attendance status of a user for a specific event
// It takes userID and eventID as input and returns the attendance status
func GetEventAttendance(userID, eventID int) (string, error) {
	var status string
	err := db.QueryRow(`
		SELECT response
		FROM Events_Responses
		WHERE user_id = ? AND event_id = ?`, userID, eventID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // No response found
		}
		log.Println("Error retrieving event attendance status:", err)
		return "", err
	}

	return status, nil
}
