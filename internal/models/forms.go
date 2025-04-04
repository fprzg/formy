package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type Form struct {
	ID           int    `json:"id"`
	UserID       int    `json:"user_id"`
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	CreatedAt    string `json:"created_at"`
	LastModified string `json:"last_modified"`
}

type FormInstance struct {
	ID        int             `json:"id"`
	FormID    int             `json:"form_id"`
	Fields    json.RawMessage `json:"fields"` // JSON string
	CreatedAt string          `json:"created_at"`
}

type FormsModelInterface interface {
	Insert(userID, name, description string) error
	GetForm(id string) (Form, error)

	InsertFormInstance(formID int, fields string) error
	GetFormInstance(formId int) (FormInstance, error)
}

type FormsModel struct {
	db *sql.DB
}

// Insert Form
func (m *FormsModel) Insert(userID, name, description string) error {
	query := `
        INSERT INTO forms (user_id, name, description)
        VALUES (?, ?, ?)
        RETURNING id, created_at, last_modified
    `

	var f Form
	err := m.db.QueryRow(query, userID, name, description).Scan(&f.ID, &f.CreatedAt, &f.LastModified)
	return err
}

func (m *FormsModel) GetForm(id string) (Form, error) {
	query := `
        SELECT id, user_id, name, description, created_at, last_modified
        FROM forms
        WHERE id = ?
    `

	var f Form
	err := m.db.QueryRow(query, id).Scan(&f.ID, &f.UserID, &f.Name, &f.Description, &f.CreatedAt, &f.LastModified)
	if err == sql.ErrNoRows {
		return Form{}, fmt.Errorf("form not found")
	}

	return f, err
}

func (m *FormsModel) InsertFormInstance(formID int, fields string) error {
	query := `
        INSERT INTO form_instances (form_id, fields)
        VALUES (?, ?)
        RETURNING id, created_at
    `

	var fi FormInstance
	err := m.db.QueryRow(query, formID, fields).Scan(&fi.ID, &fi.CreatedAt)
	return err

}

func (m *FormsModel) GetFormInstance(formId int) (FormInstance, error) {
	query := `
        SELECT id, form_id, fields, created_at
        FROM form_instances
        WHERE id = ?
    `

	var fi FormInstance
	err := m.db.QueryRow(query, formId).Scan(&fi.ID, &fi.FormID, &fi.Fields, &fi.CreatedAt)
	if err == sql.ErrNoRows {
		return FormInstance{}, fmt.Errorf("form instance not found")
	}

	return fi, err
}
