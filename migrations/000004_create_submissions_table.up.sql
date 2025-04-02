-- Up migration

CREATE TABLE submissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    form_id INTEGER NOT NULL,
    data TEXT NOT NULL, -- JSON string storing the form contents
    submitted_at DATETIME DEFAULT ( datetime('now', 'utc')),
    ip_address TEXT,
    user_agent TEXT,
    status TEXT DEFAULT 'pending', -- 'pending', 'processed', 'error'
    FOREIGN KEY (form_id) REFERENCES forms(id) ON DELETE CASCADE
);

CREATE INDEX idx_submissions_form_id ON submissions(form_id);
CREATE INDEX idx_submissions_submitted_at ON submissions(submitted_at);