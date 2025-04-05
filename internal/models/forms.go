package models

import (
	"database/sql"
	"fmt"
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
	UpdateForm(formID int, name, description, fields string) error
	DeleteForm(formID int) error
}

type FormsModel struct {
	db *sql.DB
}

func (m *FormsModel) InsertForm(userID int, name, description, fields string) error {
	const queryForm = `
        INSERT INTO forms (user_id, name, description)
        VALUES (?, ?, ?)
        RETURNING id, created_at, last_modified
    `

	var f Form
	err := m.db.QueryRow(queryForm, userID, name, description).Scan(&f.ID, &f.CreatedAt, &f.LastModified)
	if err != nil {
		return err
	}

	// TODO(Farid): Update the last_modified field on the referenced form

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
        SELECT id, name, description, created_at, last_modified
        FROM forms
        WHERE id = ?
    `

	var f Form
	err := m.db.QueryRow(query, formID).Scan(&f.ID, &f.Name, &f.Description, &f.CreatedAt, &f.LastModified)
	if err == sql.ErrNoRows {
		return Form{}, fmt.Errorf("form not found")
	}

	return f, err
}

func (m *FormsModel) GetFormsByUser(userID int) ([]Form, error) {
	const query = `
    SELECT id, name, description, created_at, last_modified
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

func (m *FormsModel) UpdateForm(id int, name, description, fields string) error {
	return nil
}

func (m *FormsModel) DeleteForm(id int) error {
	return nil
}
