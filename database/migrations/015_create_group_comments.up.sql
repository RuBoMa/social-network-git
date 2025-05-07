CREATE TABLE IF NOT EXISTS group_comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT,
    image_path TEXT,
    created_at DATETIME,
    updated_at DATETIME,
    FOREIGN KEY (group_post_id) REFERENCES group_posts(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);