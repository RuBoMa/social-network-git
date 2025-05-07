CREATE TABLE IF NOT EXISTS group_members (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    is_admin BOOLEAN,
    joined_at DATETIME,
    FOREIGN KEY (group_id) REFERENCES groups_table(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);