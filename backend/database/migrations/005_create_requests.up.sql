CREATE TABLE requests (
    id INTEGER PRIMARY KEY,
    sent_id INTEGER,
    received_id INTEGER,
    group_id INTEGER,
    status TEXT,
    created_at DATETIME,
    updated_at DATETIME,
    FOREIGN KEY (sent_id) REFERENCES users(id),
    FOREIGN KEY (received_id) REFERENCES users(id),
    FOREIGN KEY (group_id) REFERENCES groups_table(id)
);