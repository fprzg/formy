CREATE TABLE forms (
    id TEXT PRIMARY KEY,
    user_id INTEGER,
    name TEXT,
    created_at INTEGER
    FOREIGN KEY (user_id) REFERENCES users(id);
);