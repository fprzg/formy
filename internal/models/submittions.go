package models

import (
	"context"
	"database/sql"

	"formy.fprzg.net/internal/types"
	"github.com/labstack/echo/v4"
)

type SubmissionsModelInterface interface {
	Insert(submission types.SubmissionData, ctx context.Context) (int, error)
	//GetData(submissionID int) (types.SubmissionData, error)
	CheckForRepeatedUniqueField(formInstanceID int, fieldName, fieldHash string) (bool, error)
}

type SubmissionsModel struct {
	db *sql.DB
	e  *echo.Echo
}

func (m *SubmissionsModel) Insert(submission types.SubmissionData, ctx context.Context) (int, error) {
	m.e.Logger.Printf("Insert: starting submission insert for form ID %d.\n", submission.FormID)

	ctx, cancel := context.WithTimeout(ctx, contextDuration)
	defer cancel()

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		m.e.Logger.Printf("Insert: failed to begin transaction: '%v'.\n", err)
		return 0, err
	}
	defer func() {
		if err != nil {
			m.e.Logger.Printf("Insert: rolling back transaction due to error: '%v'.\n", err)
			tx.Rollback()
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				m.e.Logger.Printf("Insert: commit failed: '%v'.\n", commitErr)
				err = commitErr
			}
		}
	}()

	const stmtSubmissionsInsert = `
		INSERT INTO submissions (form_id, form_instance_id, metadata)
		VALUES (?, ?, ?)
		RETURNING id, submitted_at
	`
	err = tx.QueryRowContext(ctx, stmtSubmissionsInsert, submission.FormID, submission.FormInstanceID, submission.Metadata).Scan(&submission.ID, &submission.SubmittedAt)
	if err != nil {
		m.e.Logger.Printf("Insert: failed to insert submission: '%v'.\n", err)
		return 0, err
	}

	for _, field := range submission.Fields {
		const stmt = `
			INSERT INTO submission_fields (submission_id, field_name, content)
			VALUES (?, ?, ?)
		`
		_, err = tx.ExecContext(ctx, stmt, submission.ID, field.Name, field.Content)
		if err != nil {
			m.e.Logger.Printf("Insert: failed to insert field '%s': '%v'.\n", field.Name, err)
			return 0, err
		}

		if field.Unique {
			const stmt = `
				INSERT INTO unique_submission_fields (submission_id, form_instance_id, field_name, field_hash)
				VALUES (?, ?, ?, ?)
			`
			_, err = tx.ExecContext(ctx, stmt, submission.ID, submission.FormInstanceID, field.Name, field.Hash)
			if err != nil {
				m.e.Logger.Printf("Insert: failed to insert unique field '%s': '%v'.\n", field.Name, err)
				return 0, err
			}
		}
	}

	m.e.Logger.Printf("Insert: submission inserted successfully with ID %d.\n", submission.ID)
	return submission.ID, nil
}

func (m *SubmissionsModel) GetData(submissionID int) (types.SubmissionData, error) {
	return types.SubmissionData{}, nil
}

func (m *SubmissionsModel) CheckForRepeatedUniqueField(formInstanceID int, fieldName, fieldHash string) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT form_instance_id, field_name, field_hash
			FROM unique_submission_fields
			WHERE form_instance_id = ?
				AND field_name = ?
				AND field_hash = ?
		)
	`

	var exists bool
	err := m.db.QueryRow(query, formInstanceID, fieldName, fieldHash).Scan(&exists)
	return exists, err
}
