-- Up migration

PRAGMA foreign_keys = ON;

CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_login DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT 0,
    token_version INTEGER DEFAULT 0
);

CREATE INDEX idx_users_email ON users(email);

--
--      EXAMPLE
--
-- INSERT INTO users (name, email, password)
-- VALUES ('John Doe', 'john@doe.com', 'hashed_password');
---


CREATE TABLE refresh_tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    token TEXT NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
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
---