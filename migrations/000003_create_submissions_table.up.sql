-- Up migration

CREATE TABLE submissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    form_id INTEGER NOT NULL,
    form_instance_id INTEGER NOT NULL,
    metadata TEXT NOT NULL CHECK (json_valid(metadata)),
    submitted_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (form_id) REFERENCES forms(id),
    FOREIGN KEY (form_instance_id) REFERENCES form_instances(id)
);

CREATE INDEX idx_submissions_form_instance_id ON submissions(form_instance_id);

CREATE INDEX idx_submissions_form_instance_id_submitted_at
ON submissions(form_instance_id, submitted_at);

CREATE TABLE submission_fields (
    field_name TEXT NOT NULL,
    submission_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    PRIMARY KEY (field_name, submission_id),
    FOREIGN KEY (submission_id) REFERENCES submissions(id)
);

CREATE INDEX idx_submission_fields_submission_id ON submission_fields(submission_id);

CREATE TABLE unique_submission_fields (
    form_instance_id INTEGER NOT NULL,
    field_name TEXT NOT NULL,
    field_hash TEXT NOT NULL,
    submission_id INTEGER NOT NULL,
    PRIMARY KEY (form_instance_id, field_name, field_hash),
    FOREIGN KEY (form_instance_id) REFERENCES form_instances(id),
    FOREIGN KEY (submission_id) REFERENCES submissions(id)
);

CREATE INDEX idx_unique_submission_fields_submission_id
ON unique_submission_fields(submission_id);
