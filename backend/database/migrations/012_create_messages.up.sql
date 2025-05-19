CREATE TABLE IF NOT EXISTS Messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sender_id INTEGER,
    received_id INTEGER,
    group_id INTEGER,
    content TEXT,
    is_read BOOLEAN,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_id) REFERENCES Users(id),
    FOREIGN KEY (received_id) REFERENCES Users(id),
    FOREIGN KEY (group_id) REFERENCES Groups_Table(id)
);