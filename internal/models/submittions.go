package models

import (
	"context"
	"database/sql"
	"log"
	"time"

	"formy.fprzg.net/internal/types"
)

type SubmissionsModelInterface interface {
	Insert(submission types.SubmissionData, ctx context.Context) (int, error)
	GetData(submissionID int) (types.SubmissionData, error)
	CheckForRepeatedUniqueField(formInstanceID int, fieldName, fieldHash string) (bool, error)
	TransactionDuration() time.Duration
}

type SubmissionsModel struct {
	db                  *sql.DB
	transactionDuration time.Duration
}

func (m *SubmissionsModel) Insert(submission types.SubmissionData, ctx context.Context) (int, error) {
	log.Printf("Insert: starting submission insert for form ID %d", submission.FormID)

	ctx, cancel := context.WithTimeout(ctx, m.transactionDuration)
	defer cancel()

	// BEGIN TRANSACTION
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Insert: failed to begin transaction: %v", err)
		return 0, err
	}
	defer func() {
		if err != nil {
			log.Printf("Insert: rolling back transaction due to error: %v", err)
			tx.Rollback()
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				log.Printf("Insert: commit failed: %v", commitErr)
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
		log.Printf("Insert: failed to insert submission: %v", err)
		return 0, err
	}

	for _, field := range submission.Fields {
		const stmt = `
			INSERT INTO submission_fields (submission_id, field_name, content)
			VALUES (?, ?, ?)
		`
		_, err = tx.ExecContext(ctx, stmt, submission.ID, field.Name, field.Content)
		if err != nil {
			log.Printf("Insert: failed to insert field %s: %v", field.Name, err)
			return 0, err
		}

		if field.Unique {
			const stmt = `
				INSERT INTO unique_submission_fields (submission_id, form_instance_id, field_name, field_hash)
				VALUES (?, ?, ?, ?)
			`
			_, err = tx.ExecContext(ctx, stmt, submission.ID, submission.FormInstanceID, field.Name, field.Hash)
			if err != nil {
				log.Printf("Insert: failed to insert unique field %s: %v", field.Name, err)
				return 0, err
			}
		}
	}

	log.Printf("Insert: submission inserted successfully with ID %d", submission.ID)
	return submission.ID, nil
}

func (m *SubmissionsModel) GetData(id int) (types.SubmissionData, error) {
	/*
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
	*/
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

func (m *SubmissionsModel) TransactionDuration() time.Duration {
	return m.transactionDuration
}

/*
func (m *SubmissionsModel) getFields(submissionID int) ([]SubmissionField, error) {
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
		if err := rows.Scan(&sf.FieldName, &sf.SubmissionID, &sf.FieldContent); err != nil {
			//sendError(w, http.StatusInternalServerError, "Failed to scan submission fields: "+err.Error())
			return nil, fmt.Errorf("failed to scan submission fields: " + err.Error())
		}
		fields = append(fields, sf)
	}

	return fields, nil
}

func (m *SubmissionsModel) getUniqueSubmissionFields(formInstanceID int) ([]UniqueSubmissionField, error) {
	query := `
        SELECT field_name, field_hash, submission_id
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
		if err := rows.Scan(&usf.FieldName, &usf.FieldHash, &usf.SubmissionID); err != nil {
			return nil, fmt.Errorf("failed to scan unique submission fields: " + err.Error())
		}
		fields = append(fields, usf)
	}

	return fields, nil
}
*/
