package seed

import (
	"database/sql"
	"time"
)

type Group struct {
	ID          int
	CreatorID   int
	Title       string
	Description string
}

func SeedGroups(db *sql.DB) error {
	// Define the SQL statement to insert a group
	groups := []Group{
		{ID: 1, CreatorID: 1, Title: "Developers", Description: "A group for cool developers"},
	}

	for _, group := range groups {
		_, err := db.Exec(`INSERT OR IGNORE INTO Groups_table 
		(id, creator_id, title, description, created_at)
		VALUES (?, ?, ?, ?, ?)`,
			group.ID, group.CreatorID, group.Title, group.Description, time.Now())
		if err != nil {
			return err
		}
	}

	// Seed group members
	if err := SeedGroupMembers(db, groups); err != nil {
		return err
	}

	return nil
}

func SeedGroupMembers(db *sql.DB, groups []Group) error {
	// Define the SQL statement to insert group members
	for _, group := range groups {
		if group.ID == 0 {
			return nil // Skip if group ID is not set
		}
		groupMembers := []struct {
			ID      int
			GroupID int
			UserID  int
			IsAdmin bool
		}{
			{ID: 1, GroupID: group.ID, UserID: 1, IsAdmin: true},
			{ID: 2, GroupID: group.ID, UserID: 2, IsAdmin: false},
		}

		for _, member := range groupMembers {
			_, err := db.Exec(`
			INSERT OR IGNORE INTO Group_members 
			(id, group_id, user_id, is_admin, joined_at) 
			VALUES (?, ?, ?, ?, ?)`,
				member.ID, member.GroupID, member.UserID, member.IsAdmin, time.Now())
			if err != nil {
				return err
			}
		}
	}

	return nil
}
