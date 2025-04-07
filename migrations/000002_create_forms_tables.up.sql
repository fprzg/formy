-- Up migration

PRAGMA foreign_keys = ON;

CREATE TABLE forms (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    form_version INTEGER DEFAULT 1,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_forms_user_id ON forms(user_id);

CREATE TABLE form_instances (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    form_id INTEGER NOT NULL,
    form_version INTEGER NOT NULL,
    fields TEXT NOT NULL CHECK (json_valid(fields)),
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (form_id) REFERENCES forms(id) ON DELETE CASCADE
);