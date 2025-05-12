package models

type Post struct {
	PostID      int       `json:"post_id"`
	Author      User      `json:"author"`
	Group       Group     `json:"group"`
	PostTitle   string    `json:"post_title"`
	PostContent string    `json:"post_content"`
	Comments    []Comment `json:"comments"`
	PostImage   string    `json:"post_image"`
	Privacy     string    `json:"privacy"`
	CreatedAt   string    `json:"created_at"`
}
type Comment struct {
	CommentID      int    `json:"comment_id"`
	PostID         int    `json:"post_id"`
	CommentAuthor  User   `json:"comment_author"`
	CommentContent string `json:"comment_content"`
	CommentImage   string `json:"comment_image"`
	CreatedAt      string `json:"created_at"`
}

// User details
type User struct {
	UserID      int    `json:"user_id"`
	Nickname    string `json:"nickname"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DateOfBirth string `json:"date_of_birth"`
	Email       string `json:"email"`
	AboutMe     string `json:"about_me"`
	AvatarPath  string `json:"avatar_path"`
	IsPublic    bool   `json:"is_public"`
	Password    string `json:"password"`
}

// ProfileResponse contains all data needed for a profile page
type ProfileResponse struct {
	User           User   `json:"user"`
	IsOwnProfile   bool   `json:"is_own_profile"`
	IsFollower     bool   `json:"is_follower"`
	Posts          []Post `json:"posts"`
	FollowersCount int    `json:"followers_count"`
	FollowingCount int    `json:"following_count"`
}

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

// Follow request
type Request struct {
	RequestID int    `json:"request_id"`
	Sender    User   `json:"sender"`
	Receiver  User   `json:"receiver"`
	Group     Group  `json:"group"`
	Status    string `json:"status"` // "pending", "accepted", "rejected"
	CreatedAt string `json:"created_at"`
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

// Notifications
type Notification struct {
	NotificationID int     `json:"notification_id"`
	UserID         int     `json:"user_id"`
	Type           string  `json:"type"`
	Event          Event   `json:"event"`
	FollowRequest  Request `json:"follow_request"`
	Seen           bool    `json:"seen"`
	CreatedAt      string  `json:"created_at"`
}

// Struct to map the incoming login data
type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Struct to map json response messages NEEDS MORE THAN A MESSAGE ??
type Response struct {
	Message string `json:"message"`
}
