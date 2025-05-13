CREATE TABLE IF NOT EXISTS Group_Invitations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER,
    sender_id INTEGER,
    reciver_id INTEGER,
    status TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    FOREIGN KEY (group_id) REFERENCES Groups(id),
    FOREIGN KEY (sender_id) REFERENCES Users(id),
    FOREIGN KEY (reciver_id) REFERENCES Users(id)
);