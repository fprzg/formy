-- Up migration

CREATE TABLE refresh_tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    token TEXT NOT NULL UNIQUE,
    created_at DATETIME DEFAULT (datetime('now', 'utc')),
    expires_at DATETIME NOT NULL,
    revoked BOOLEAN DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);

--
--      EXAMPLE
--
-- INSERT INTO refresh_tokens (user_id, token, expires_at)
-- VALUES (1, 'some_token', datetime('2025-04-02 12:00:00', '+6 hours', 'utc'));