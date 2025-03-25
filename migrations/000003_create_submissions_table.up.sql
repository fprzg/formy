CREATE TABLE submissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    form_id INTEGER NOT NULL,
    enforced_field TEXT,
    data TEXT,
    FOREIGN KEY (form_id) REFERENCES forms(form_id);
);