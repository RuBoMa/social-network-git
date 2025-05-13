CREATE TABLE IF NOT EXISTS Follow_Requests (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sender_id INTEGER,
    reciver_id INTEGER,
    status TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    FOREIGN KEY (sender_id) REFERENCES Users(id),
    FOREIGN KEY (reciver_id) REFERENCES Users(id)
);