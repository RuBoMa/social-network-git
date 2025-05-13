package models

// Group details
type Group struct {
	GroupID        int    `json:"group_id"`
	GroupName      string `json:"group_name"`
	GroupDesc      string `json:"group_desc"`
	GroupCreator   User   `json:"group_creator"`
	GroupCreatedAt string `json:"group_created_at"`
	GroupMembers   []User `json:"group_members"`
	GroupPosts     []Post `json:"group_posts"`
}

// Group members
type GroupMember struct {
	GroupID int    `json:"group_id"`
	Members []User `json:"members"`
}

// Event details
type Event struct {
	EventID     int    `json:"event_id"`
	CreatorID   int    `json:"creator_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	EventDate   string `json:"event_date"`
	CreatedAt   string `json:"created_at"`
	Group       Group  `json:"group"`
}
