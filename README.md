# real-time-forum

This is a web-based forum that allows registered users to communicate, post content, interact with each other's posts and comments, and send private messages to online users. <br>

### Project structure

```
forum/
├── main.go                    # Entry point of the application
├── index.html                 # SPA used to render all pages
├── go.mod                     # Go module definition
├── go.sum                     # Dependency checksums
├── .gitignore                 # Files and directories to ignore in Git
├── README.md                  # Project documentation
├── database/
│   ├── database.go            # Handles database connection setup
│   ├── queries.go             # Contains reusable database query logic
│   └── table.go               # Includes table creation or schema management logic
├── js/
│   ├── app.js                 # Handles which page is loaded, authentication and API calls
│   ├── chat.js                # Handles chat functionality (sending and loading messages)
│   ├── createpost.js          # Injects HTML with content for creating a post
│   ├── error.js               # Displays error page
│   ├── feed.js                # Injects HTML with all posts
│   ├── login.js               # Handles the login form
│   ├── logout.js              # Logs the user out by clearing the session
│   ├── post.js                # Injects HTML with a specific posts, handles comment and likes
│   ├── session.js             # Refreshed the session every 15 min
│   └── sidebar.js             # Handles the websocket messages and updating the sidebar
│   └── signup.js              # Handles the signup form
├── assets/
│   ├── favicon.ico            # Favicon for the application
│   └── styles.css             # CSS file for frontend styling
├── backend/
    ├── APIHandler.go          # API call handling
    ├── createpost.go          # Handles logic for creating posts
    ├── dbentries.go           # Reusable functions to make database entries
    ├── dbquery.go             # Reusable functions for make database queries
    ├── feed.go                # Root page logic
    ├── login.go               # Handles user login
    ├── logout.go              # Handles user logout
    ├── post.go                # Handles post viewing, commenting and likes
    ├── response.go            # Sends simple responses
    ├── session.go             # Authenticates, creates and expires sessions
    ├── signup.go              # Handles user registration
    └── structs.go             # Shared data structures
```

- **Backend:** The backend is built using Golang, which handles data manipulation and WebSocket connections.

- **Frontend:** The frontend is developed using JavaScript, HTML, and CSS, which manage user interactions, page updates, and client-side WebSocket connections.

### Authentication

Once the user logs in, they are given a UUID token in a session cookie. This token is used to authenticate the user.<br>
The token is stored in the database and is used to authenticate the user. When an user logs out, the token is marked as deleted in the database.<br>

## Features

- **User Registration & Authentication**:

  - Users can register by providing their email, username, age, gender, first name, last name and password.
  - Passwords are encrypted.
  - Only registered users can access the site.
  - Live chat between online users is manages via Websocket.
  - Sessions are managed using cookies.

- **SQLite Database**:

  - Data such as users, posts, comments, and likes/dislikes are stored in an SQLite database.
  - At least one SELECT, CREATE, and INSERT SQL query is used.

## Setup & Installation

1. **Clone the repository**:
   ```bash
   git clone https://01.gritlab.ax/git/fvesanen/real-time-forum
   cd real-time-forum
   ```
2. **Install dependencies:**
   ```bash
   go get github.com/gorilla/websocket
   go get github.com/mattn/go-sqlite3
   go get golang.org/x/crypto/bcrypt
   go get github.com/gofrs/uuid
   ```
3. Run the Go server:
   ```bash
   go run main.go
   ```

### Technologies Used

**Go:** Backend server and logic.<br>
**SQLite:** Local database for storing users, posts, comments, likes/dislikes.<br>
**Cookies & Sessions:** For user authentication and maintaining logged-in sessions.<br>
**HTML:** Key element of the Single Page Application, one HTML is used to display different pages.<br>
**JavaScript:** Handles frontend interactions, dynamic content updates, and client-side WebSocket connections for real-time communication.<br>
**WebSockets:** For real-time message notifications, user list updates, typing in progress and private messaging between users.<br>
