package models

// Group details
type Group struct {
	GroupID        int       `json:"group_id"`
	GroupName      string    `json:"group_name"`
	GroupDesc      string    `json:"group_desc"`
	GroupCreator   User      `json:"group_creator"`
	GroupCreatedAt string    `json:"group_created_at"`
	GroupMembers   []User    `json:"group_members"`
	GroupPosts     []Post    `json:"group_posts"`
	GroupEvents    []Event   `json:"group_events"`
	GroupRequests  []Request `json:"group_requests"`
	IsMember       bool      `json:"is_member"`      // true if the user is a member of the group
	RequestStatus  string    `json:"request_status"` // "requested", "invited"
	RequestID      int       `json:"request_id"`     // ID of the request if it exists
}

// Group members
type GroupMember struct {
	GroupID int    `json:"group_id"`
	Members []User `json:"members"`
}

// Event details
type Event struct {
	EventID        int             `json:"event_id"`
	CreatorID      int             `json:"creator_id"`
	Title          string          `json:"title"`
	Description    string          `json:"description"`
	EventDate      string          `json:"event_date"`
	CreatedAt      string          `json:"created_at"`
	Group          Group           `json:"group"`
	EventResponses []EventResponse `json:"event_responses"`
}

type EventResponse struct {
	ResponseID int    `json:"response_id"`
	Event      Event  `json:"event"`
	User       User   `json:"user"`
	Response   string `json:"response"` // "going", "not going"
}
