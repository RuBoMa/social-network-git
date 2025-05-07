CREATE TABLE IF NOT EXISTS group_join_requests (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    status TEXT,
    created_at DATETIME,
    updated_at DATETIME,
    FOREIGN KEY (group_id) REFERENCES groups_table(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);