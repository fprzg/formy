package models

import (
	"database/sql"
	"strings"
)

type Form struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	CreatedAt    string `json:"created_at"`
	LastModified string `json:"last_modified"`
}

type formData struct {
	UserID int `json:"user_id"`
	Form
}

type FormInstance struct {
	ID        int    `json:"id"`
	FormID    int    `json:"form_id"`
	Fields    string `json:"fields"` // JSON string
	CreatedAt string `json:"created_at"`
}

type FormsModelInterface interface {
	InsertForm(userID int, name, description, fields string) error
	GetForm(formID int) (Form, error)
	GetFormsByUser(userID int) ([]Form, error)
	UpdateFormName(formID int, name string) error
	UpdateFormDescription(formID int, description string) error
	DeleteForm(formID int) error
}

type FormsModel struct {
	db *sql.DB
}

func (m *FormsModel) InsertForm(userID int, name, description, fields string) error {
	if userID < 1 {
		return ErrInvalidUserID
	}
	if name == "" || fields == "" {
		return ErrInvalidInput
	}
	const queryForm = `
        INSERT INTO forms (user_id, name, description)
        VALUES (?, ?, ?)
        RETURNING id, created_at, updated_at
    `

	var f Form
	err := m.db.QueryRow(queryForm, userID, name, description).Scan(&f.ID, &f.CreatedAt, &f.LastModified)
	if err != nil {
		if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
			return ErrInvalidUserID
		}
		return err
	}

	const queryFormInstance = `
		INSERT INTO form_instances (form_id, fields)
		VALUES (?, ?)
		RETURNING id, created_at
	`

	var fi FormInstance
	err = m.db.QueryRow(queryFormInstance, f.ID, fields).Scan(&fi.ID, &fi.CreatedAt)

	return err
}

func (m *FormsModel) GetForm(formID int) (Form, error) {
	//SELECT id, user_id, name, description, created_at, last_modified
	const query = `
        SELECT id, name, description, created_at, updated_at
        FROM forms
        WHERE id = ?
    `

	var f Form
	err := m.db.QueryRow(query, formID).Scan(&f.ID, &f.Name, &f.Description, &f.CreatedAt, &f.LastModified)
	if err == sql.ErrNoRows {
		return Form{}, ErrFormNotFound
	}

	return f, err
}

func (m *FormsModel) GetFormsByUser(userID int) ([]Form, error) {
	const query = `
    SELECT id, name, description, created_at, updated_at
	FROM forms
	WHERE user_id = ?
	`

	rows, err := m.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var forms []Form
	for rows.Next() {
		var f Form
		err = rows.Scan(&f.ID, &f.Name, &f.Description, &f.CreatedAt, &f.LastModified)
		if err != nil {
			return nil, err
		}

		forms = append(forms, f)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return forms, nil
}

func (m *FormsModel) UpdateFormName(formID int, name string) error {
	const stmt = `
	UPDATE forms
	SET name = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?
	`

	rows, err := ExecuteSqlStmt(m.db, stmt, name, formID)
	if rows == 0 {
		return ErrFormNotFound
	}

	return err
}

func (m *FormsModel) UpdateFormDescription(formID int, description string) error {
	const query = `
	UPDATE forms
	SET description = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?
	`

	rows, err := ExecuteSqlStmt(m.db, query, description, formID)
	if rows == 0 {
		return ErrFormNotFound
	}

	return err
}

func (m *FormsModel) DeleteForm(formID int) error {
	const stmt = `
	DELETE FROM forms WHERE id = ?
	`

	rows, err := ExecuteSqlStmt(m.db, stmt, formID)
	if rows == 0 {
		return ErrFormNotFound
	}

	return err
}
