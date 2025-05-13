CREATE TABLE IF NOT EXISTS Events_Responses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id INTEGER,
    user_id INTEGER,
    response TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    FOREIGN KEY (event_id) REFERENCES Events(id),
    FOREIGN KEY (user_id) REFERENCES Users(id)
);