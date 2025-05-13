CREATE TABLE IF NOT EXISTS Group_Members (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER,
    user_id INTEGER,
    is_admin BOOLEAN,
    joined_at DATETIME,
    FOREIGN KEY (group_id) REFERENCES Groups_Table(id),
    FOREIGN KEY (user_id) REFERENCES Users(id)
);