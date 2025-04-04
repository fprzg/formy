package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type SubmissionField struct {
	FieldName    string `json:"field_name"`
	SubmissionID int    `json:"submission_id"`
	Content      string `json:"content"`
}

type Submission struct {
	ID             int             `json:"id"`
	FormID         int             `json:"form_id"`
	FormInstanceID int             `json:"form_instance_id"`
	Metadata       json.RawMessage `json:"metadata"` // JSON string
	SubmittedAt    string          `json:"submitted_at"`
}

type SubmissionsModelInterface interface {
}

type SubmissionsModel struct {
	db *sql.DB
}

func (m *SubmissionsModel) InsertSubmission(formID, formInstanceID int, metadata string) error {
	query := `
        INSERT INTO submissions (form_id, form_instance_id, metadata)
        VALUES (?, ?, ?)
        RETURNING id, submitted_at
    `

	var s Submission
	err := m.db.QueryRow(query, formID, formInstanceID, metadata).Scan(&s.ID, &s.SubmittedAt)
	return err
}

func (m *SubmissionsModel) GetSubmission(id int) (Submission, error) {
	query := `
        SELECT id, form_id, form_instance_id, metadata, submitted_at
        FROM submissions
        WHERE id = ?
    `
	var s Submission
	err := m.db.QueryRow(query, id).Scan(&s.ID, &s.FormID, &s.FormInstanceID, &s.Metadata, &s.SubmittedAt)
	if err == sql.ErrNoRows {
		return Submission{}, fmt.Errorf("submission not found")
	}

	return s, err
}

func (m *SubmissionsModel) InsertSubmissionField(submissionID, fieldName, content string) error {
	query := `
        INSERT INTO submission_fields (field_name, submission_id, content)
        VALUES (?, ?, ?)
    `

	_, err := m.db.Exec(query, fieldName, submissionID, content)

	return err
}

func (m *SubmissionsModel) GetSubmissionFields(submissionID int) ([]SubmissionField, error) {
	query := `
        SELECT field_name, submission_id, content
        FROM submission_fields
        WHERE submission_id = ?
    `

	rows, err := m.db.Query(query, submissionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get submission fields: " + err.Error())
	}
	defer rows.Close()

	var fields []SubmissionField
	for rows.Next() {
		var sf SubmissionField
		if err := rows.Scan(&sf.FieldName, &sf.SubmissionID, &sf.Content); err != nil {
			//sendError(w, http.StatusInternalServerError, "Failed to scan submission fields: "+err.Error())
			return nil, fmt.Errorf("failed to scan submission fields: " + err.Error())
		}
		fields = append(fields, sf)
	}

	return fields, nil
}

type UniqueSubmissionField struct {
	FormInstanceID int    `json:"form_instance_id"`
	FieldName      string `json:"field_name"`
	FieldHash      string `json:"field_hash"`
	SubmissionID   int    `json:"submission_id"`
}

func (m *SubmissionsModel) InsertUniqueSubmissionField(formInstanceID, submissionID int, fieldName string, fieldHash []byte) error {
	query := `
        INSERT INTO unique_submission_fields (form_instance_id, field_name, field_hash, submission_id)
        VALUES (?, ?, ?, ?)
    `

	_, err := m.db.Exec(query, formInstanceID, fieldName, fieldHash, submissionID)
	if err != nil {
		return fmt.Errorf("failed to insert unique submission field: " + err.Error())
	}
	return nil
}

func (m *SubmissionsModel) GetUniqueSubmissionFields(formInstanceID int) ([]UniqueSubmissionField, error) {
	query := `
        SELECT form_instance_id, field_name, field_hash, submission_id
        FROM unique_submission_fields
        WHERE form_instance_id = ?
    `

	rows, err := m.db.Query(query, formInstanceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get unique submission fields" + err.Error())
	}
	defer rows.Close()

	var fields []UniqueSubmissionField
	for rows.Next() {
		var usf UniqueSubmissionField
		if err := rows.Scan(&usf.FormInstanceID, &usf.FieldName, &usf.FieldHash, &usf.SubmissionID); err != nil {
			return nil, fmt.Errorf("failed to scan unique submission fields: " + err.Error())
		}
		fields = append(fields, usf)
	}

	return fields, nil
}
