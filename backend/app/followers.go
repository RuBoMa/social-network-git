package app

// When someone has pressed "follow"
func HandleNewFollower() {

	// check if profile public or private

	// IF PUBLIC
	// Check if user is already in the follow table
	// Change status / add user to followers table

	// IF PRIVATE
	// Check if user ir already in the follow_requests table
	// Change status / add user to the table as pending

	// Send notification about new follower / follow request

}

// Accept or decline to follow request
func HandleFollowRequest() {

	// IF ACCEPT
	// Change status to accept
	// Change status / add user to Followers table
	// Send notification about acceptance status

	// IF DECLINE
	// Change status to decline

}

// When user pressed "unfollow" on someone
func HandleUnfollow() {

	// Change status on followers table
	// Check if the user is a post_privacy table
	// if yes, change status to deleted

}
