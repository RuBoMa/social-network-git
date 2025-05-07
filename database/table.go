package database

import (
	"database/sql"
	"fmt"
)

// MakeTables creates the tables in the database if they do not exist and inserts initial data into the tables
func MakeTables(db *sql.DB) {

	createUserTableQuery := `
		CREATE TABLE IF NOT EXISTS User (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL COLLATE NOCASE,
		age INTEGER NOT NULL,
		gender TEXT NOT NULL CHECK (gender IN ('male', 'female', 'non-binary', 'other', 'prefer not to say')),
		firstname TEXT NOT NULL,
		lastname TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at TEXT NOT NULL,
		updated_at TEXT,
		status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'deleted'))

	);`
	if _, err := db.Exec(createUserTableQuery); err != nil {
		fmt.Println("Error creating User table:", err)
		return
	}
	createPostTableQuery := `
		CREATE TABLE IF NOT EXISTS Post (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
    	content TEXT NOT NULL,
   		user_id INTEGER NOT NULL,
   		created_at TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'deleted')),
    	FOREIGN KEY (user_id) REFERENCES User (id)
	);`
	if _, err := db.Exec(createPostTableQuery); err != nil {
		fmt.Println("Error creating Post table:", err)
		return
	}
	createCommentTableQuery := `
		CREATE TABLE IF NOT EXISTS Comment (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
    	content TEXT NOT NULL,
    	post_id INTEGER NOT NULL,
    	user_id INTEGER NOT NULL,
    	created_at TEXT NOT NULL,
    	status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'deleted')),
    	FOREIGN KEY (post_id) REFERENCES Post (id),
    	FOREIGN KEY (user_id) REFERENCES User (id)
	);`
	if _, err := db.Exec(createCommentTableQuery); err != nil {
		fmt.Println("Error creating Comment table:", err)
		return
	}
	createCategoryTableQuery := `
		CREATE TABLE IF NOT EXISTS Category (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
    	name TEXT UNIQUE NOT NULL,
    	status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'deleted'))
	);`
	if _, err := db.Exec(createCategoryTableQuery); err != nil {
		fmt.Println("Error creating Category table:", err)
		return
	}
	createPost_CategoryTableQuery := `
		CREATE TABLE IF NOT EXISTS Post_Category (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	category_id INTEGER NOT NULL,
    	post_id INTEGER NOT NULL,
    	FOREIGN KEY (category_id) REFERENCES Category (id),
    	FOREIGN KEY (post_id) REFERENCES Post (id)
	);`
	if _, err := db.Exec(createPost_CategoryTableQuery); err != nil {
		fmt.Println("Error creating Post_Category table:", err)
		return
	}
	createLikeTableQuery := `
		CREATE TABLE IF NOT EXISTS Like (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
    	user_id INTEGER NOT NULL,
    	post_id INTEGER,
    	comment_id INTEGER,
   		created_at TEXT NOT NULL,
		updated_at TEXT,
    	type INTEGER NOT NULL CHECK (type IN (0, 1, 2)),
    	FOREIGN KEY (user_id) REFERENCES User (id),
    	FOREIGN KEY (post_id) REFERENCES Post (id),
    	FOREIGN KEY (comment_id) REFERENCES Comment (id)
	);`
	if _, err := db.Exec(createLikeTableQuery); err != nil {
		fmt.Println("Error creating Like table:", err)
		return
	}
	createSessionTableQuery := `
		CREATE TABLE IF NOT EXISTS Session (
		id TEXT PRIMARY KEY, -- Unique session ID (UUID)
    	user_id INTEGER NOT NULL,
    	status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'expired', 'inactive', 'deleted')),
    	created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL,
		expires_at TEXT NOT NULL,
		last_access TEXT NOT NULL,
    	FOREIGN KEY (user_id) REFERENCES User (id)
	);`
	if _, err := db.Exec(createSessionTableQuery); err != nil {
		fmt.Println("Error creating Session table:", err)
		return
	}

	createChatTableQuery := `
	CREATE TABLE IF NOT EXISTS Chat (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
   		user1_id INTEGER NOT NULL CHECK (user1_id != 0),
    	user2_id INTEGER NOT NULL CHECK (user2_id != 0),
		created_at TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'deleted', 'archived')),
		FOREIGN KEY (user1_id) REFERENCES User(id),
		FOREIGN KEY (user2_id) REFERENCES User(id)
	);`
	if _, err := db.Exec(createChatTableQuery); err != nil {
		fmt.Println("Error creating Chat table:", err)
		return
	}

	createMessageTableQuery := `
	CREATE TABLE IF NOT EXISTS Message (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		chat_id INTEGER NOT NULL,
		sender_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		created_at TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'deleted', 'edited')),
		FOREIGN KEY (chat_id) REFERENCES Chat(id),
		FOREIGN KEY (sender_id) REFERENCES User(id)
	);`
	if _, err := db.Exec(createMessageTableQuery); err != nil {
		fmt.Println("Error creating Message table:", err)
		return
	}

	// Inserting categories to database
	insertCategoryQuery := `
    INSERT INTO category (name)
    SELECT 'General' WHERE NOT EXISTS (SELECT 1 FROM category WHERE name = 'General')
    UNION ALL
    SELECT 'Tutorial' WHERE NOT EXISTS (SELECT 1 FROM category WHERE name = 'Tutorial')
    UNION ALL
    SELECT 'Question' WHERE NOT EXISTS (SELECT 1 FROM category WHERE name = 'Question')
	UNION ALL
    SELECT 'Plants' WHERE NOT EXISTS (SELECT 1 FROM category WHERE name = 'Plants')
	UNION ALL
    SELECT 'Pests' WHERE NOT EXISTS (SELECT 1 FROM category WHERE name = 'Pests')
	UNION ALL
    SELECT 'Sustainability' WHERE NOT EXISTS (SELECT 1 FROM category WHERE name = 'Sustainability');
`
	if _, err := db.Exec(insertCategoryQuery); err != nil {
		fmt.Println("Error inserting into Category table:", err)
		return
	}

	//Insert initial admin data
	insertUserQuery := `
    INSERT INTO User (username, age, gender, firstname, lastname, email, password, created_at)
    SELECT 'admin', 35, 'female', 'Fanni', 'Vesanen', 'admin@example.com', 'hashedpassword', datetime('now')
    WHERE NOT EXISTS (SELECT 1 FROM User WHERE username = 'admin');
`
	if _, err := db.Exec(insertUserQuery); err != nil {
		fmt.Println("Error inserting into User table:", err)
		return
	}

	//Insert initial data into Post
	insertPostQuery := `
    INSERT INTO post (title, content, user_id, created_at) 
    SELECT 'Welcome to the forum', 'This is the first post!', 1, datetime('now')
    WHERE NOT EXISTS (
        SELECT 1 FROM post WHERE title = 'Welcome to the forum'
    );
`
	result, err := db.Exec(insertPostQuery)
	if err != nil {
		fmt.Println("Error inserting into Post table:", err)
		return
	}

	// Retrieve the last inserted ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Error retrieving last insert ID:", err)
		return
	}

	insertPostCategoryQuery := `
    INSERT INTO Post_category (post_id, category_id)
	VALUES (?, 1);
`
	if _, err := db.Exec(insertPostCategoryQuery, int(lastInsertID)); err != nil {
		fmt.Println("Error inserting into Post_category table:", err)
		return
	}

	fmt.Println("Tables created and initial data inserted successfully.")
}
