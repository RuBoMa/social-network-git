package models

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
	Token       string `json:"token"`
}

// ProfileResponse contains all data needed for a profile page
type ProfileResponse struct {
	User           User      `json:"user"`
	IsOwnProfile   bool      `json:"is_own_profile"`
	IsFollower     bool      `json:"is_follower"`
	Posts          []Post    `json:"posts"`
	FollowersCount int       `json:"followers_count"`
	FollowingCount int       `json:"following_count"`
	FollowRequests []Request `json:"follow_requests"`
}
