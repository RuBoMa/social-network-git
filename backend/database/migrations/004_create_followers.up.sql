CREATE TABLE IF NOT EXISTS Followers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    follower_id INTEGER,
    followed_id INTEGER,
    status TEXT,
    updated_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (follower_id) REFERENCES Users(id),
    FOREIGN KEY (followed_id) REFERENCES Users(id)
);