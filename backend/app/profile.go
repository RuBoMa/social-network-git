package app

// When user clicks any profile
func ServeProfile() {

	//Check if it's their own profile (compare user_ids)
	// If yes, add that information to the response

	//Check if the profile is public or private and add to responce

	// Get profile information (minus password)
	// Get all posts by the user (inlude privacy) Only titles?
	// Get following SUM (user_id == follower_id) and followers SUM (user_id == followed_id)

}

func ChangeProfileVisibility() {

	//Toggle between public and private in Users table
}

func DisplayFollowData() {

	// Check if following or followers
	// Get all ids from Followers table according to above specification
	// Get all names and pictures from the Users table and return
}
