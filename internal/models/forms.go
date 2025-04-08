package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"formy.fprzg.net/internal/types"
	"formy.fprzg.net/internal/utils"
)

/*
type formData struct {
	UserID int `json:"user_id"`
	Form
}
*/

type FormInstance struct {
	ID          int    `json:"id"`
	FormID      int    `json:"form_id"`
	FormVersion int    `json:"form_version"`
	FieldsJSON  string `json:"fields"`
	CreatedAt   string `json:"created_at"`
}

type FormsModelInterface interface {
	Insert(userID int, name, description string, fields []types.FieldData) (int, error)
	Get(formID int) (types.FormData, error)
	GetFormsByUserID(userID int) ([]types.FormData, error)
	UpdateName(formID int, name string) error
	UpdateDescription(formID int, description string) error
	DeleteForm(formID int) error
}

type FormsModel struct {
	db *sql.DB
}

func (m *FormsModel) Insert(userID int, name, description string, fields []types.FieldData) (int, error) {
	const queryForm = `
        INSERT INTO forms (user_id, name, description)
        VALUES (?, ?, ?)
        RETURNING id, created_at, updated_at, form_version
    `

	const queryFormInstance = `
		INSERT INTO form_instances (form_id, fields, form_version)
		VALUES (?, ?, ?)
		RETURNING id, created_at
	`

	if userID < 1 {
		return 0, ErrInvalidUserID
	}
	if name == "" || fields == nil {
		return 0, ErrInvalidInput
	}

	if fields == nil {
		return 0, fmt.Errorf("models: form has to have at least one field")
	}

	fieldsJSON, err := utils.ToJSON(fields)
	if err != nil {
		return 0, err
	}

	var f types.FormData
	err = m.db.QueryRow(queryForm, userID, name, description).Scan(&f.ID, &f.CreatedAt, &f.UpdatedAt, &f.FormVersion)
	if err != nil {
		if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
			return 0, ErrInvalidUserID
		}
		return 0, err
	}

	var fi FormInstance
	err = m.db.QueryRow(queryFormInstance, f.ID, fieldsJSON, f.FormVersion).Scan(&fi.ID, &fi.CreatedAt)
	if err != nil {
		return 0, err
	}

	return fi.ID, nil
}

func (m *FormsModel) Get(formID int) (types.FormData, error) {
	//SELECT id, user_id, name, description, created_at, last_modified
	const queryGetForm = `
        SELECT user_id, id, name, description, created_at, updated_at
        FROM forms
        WHERE id = ?
    `

	const queryGetFormInstance = `
	SELECT id, form_version, fields
	FROM form_instances
	WHERE form_id = ?
	ORDER BY id DESC
	LIMIT 1
	`

	var f types.FormData
	err := m.db.QueryRow(queryGetForm, formID).Scan(&f.UserID, &f.ID, &f.Name, &f.Description, &f.CreatedAt, &f.UpdatedAt)
	if err == sql.ErrNoRows {
		return types.FormData{}, ErrFormNotFound
	}

	var fi FormInstance
	err = m.db.QueryRow(queryGetFormInstance, f.ID).Scan(&fi.ID, &fi.FormVersion, &fi.FieldsJSON)
	if err != nil {
		return types.FormData{}, ErrFormNotFound
	}

	err = json.Unmarshal([]byte(fi.FieldsJSON), &f.Fields)
	if err != nil {
		return types.FormData{}, err
	}

	return f, err
}

func (m *FormsModel) GetFormsByUserID(userID int) ([]types.FormData, error) {
	const query = `
    SELECT
		f.id, f.name, f.description, f.created_at, f.updated_at,
		fi.form_version, fi.fields
	FROM forms f
	LEFT JOIN form_instances fi ON fi.id = (
		SELECT id FROM form_instances
		WHERE form_id = f.id
		ORDER BY created_at DESC
		LIMIT 1
	)
	WHERE f.user_id = ?
	`

	rows, err := m.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var forms []types.FormData
	for rows.Next() {
		var f types.FormData
		var formFields string
		err = rows.Scan(
			&f.ID, &f.Name, &f.Description, &f.CreatedAt, &f.UpdatedAt,
			&f.FormVersion, &formFields)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(formFields), &f.Fields)
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

func (m *FormsModel) UpdateName(formID int, name string) error {
	const stmt = `
	UPDATE forms
	SET name = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?
	`

	rows, err := utils.ExecuteSqlStmt(m.db, stmt, name, formID)
	if rows == 0 {
		return ErrFormNotFound
	}

	return err
}

func (m *FormsModel) UpdateDescription(formID int, description string) error {
	const query = `
	UPDATE forms
	SET description = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?
	`

	rows, err := utils.ExecuteSqlStmt(m.db, query, description, formID)
	if rows == 0 {
		return ErrFormNotFound
	}

	return err
}

func (m *FormsModel) DeleteForm(formID int) error {
	const stmt = `
	DELETE FROM forms WHERE id = ?
	`

	rows, err := utils.ExecuteSqlStmt(m.db, stmt, formID)
	if rows == 0 {
		return ErrFormNotFound
	}

	return err
}
