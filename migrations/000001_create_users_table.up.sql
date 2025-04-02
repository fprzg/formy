-- Up migration

CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at DATETIME DEFAULT (datetime('now', 'utc')),
    updated_at DATETIME DEFAULT (datetime('now', 'utc')),
    last_login DATETIME,
    is_active BOOLEAN DEFAULT 0,
    token_version INTEGER DEFAULT 0
);

CREATE INDEX idx_users_email ON users(email);

--
--      EXAMPLE
--
-- INSERT INTO users (name, email, password)
-- VALUES ('John Doe', 'john@doe.com', 'hashed_password');