package database

import (
	"social_network/models"
	"strings"
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
	result, err := db.Exec("INSERT INTO Groups (title, description, creator_id, created_at) VALUES (?, ?, ?, ?)",
		group.GroupName, group.GroupDesc, group.GroupCreator.UserID, group.GroupCreatedAt)
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
	err := db.QueryRow("SELECT COUNT(*) FROM Groups WHERE title = ?", title).Scan(&count)
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
	err := db.QueryRow("SELECT COUNT(*) FROM Groups WHERE id = ?", groupID).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}

// AddGroupInvitationIntoDB adds a new group invitation to the database with status "pending"
// It takes a models.Request object as input and inserts it into the Group_Invitations table
func AddGroupInvitationIntoDB(invitation models.Request) (int, error) {
	result, err := db.Exec("INSERT INTO Group_Invitations (sender_id, receiver_id, group_id, status, created_at) VALUES (?, ?, ?, ?, ?)",
		invitation.Sender.UserID, invitation.Receiver.UserID, invitation.Group.GroupID, invitation.Status, invitation.CreatedAt)
	if err != nil {
		return 0, err
	}

	invitationID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(invitationID), nil
}
