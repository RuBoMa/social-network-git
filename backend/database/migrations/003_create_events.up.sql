CREATE TABLE IF NOT EXISTS Events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER,
    creator_id INTEGER,
    title TEXT,
    description TEXT,
    event_time DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    FOREIGN KEY (group_id) REFERENCES Groups(id),
    FOREIGN KEY (creator_id) REFERENCES Users(id)
);