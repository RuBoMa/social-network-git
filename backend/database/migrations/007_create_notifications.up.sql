CREATE TABLE IF NOT EXISTS Notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    is_read BOOLEAN,
    related_event_id INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES Users(id),
    FOREIGN KEY (related_event_id) REFERENCES Events(id)
    FOREIGN KEY (related_event_id) REFERENCES Requests(id)
);