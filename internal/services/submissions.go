package services

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"formy.fprzg.net/internal/models"
	"formy.fprzg.net/internal/types"
)

type SubmissionService struct {
	models *models.Models
}

type SubmissionServiceInterface interface {
	ProcessSubmissionForm(formID int, r *http.Request, ctx context.Context) (int, error)
}

func (s *SubmissionService) ProcessSubmissionForm(formID int, r *http.Request, ctx context.Context) (int, error) {
	formData, err := s.models.Forms.Get(formID)
	if err != nil {
		return 0, err
	}

	submission, err := s.GetSubmissionDataFromRequest(formData, r, ctx)
	if err != nil {
		return 0, err
	}

	return s.models.Submissions.Insert(submission, ctx)
}

func (s *SubmissionService) GetSubmissionDataFromRequest(form types.FormData, r *http.Request, ctx context.Context) (types.SubmissionData, error) {
	formInstanceID, err := s.models.Forms.GetFormInstanceID(form.ID)
	if err != nil {
		return types.SubmissionData{}, err
	}

	if err = r.ParseForm(); err != nil {
		return types.SubmissionData{}, err
	}

	submission := types.SubmissionData{
		FormID:         form.ID,
		FormInstanceID: formInstanceID,
		Metadata:       `{ "user_agent": "curl uwu", "ip_address": "0.0.0.0" } `,
	}

	for fieldName, fieldContents := range r.Form {
		fieldIndex := form.GetFieldIndex(fieldName)
		if fieldIndex == -1 {
			// TODO: Report incident
			log.Printf("Insert: unknown field %s; skipping", fieldName)
			continue
		}

		formField := form.Fields[fieldIndex]

		subField := types.SubmissionField{
			Name:    formField.Name,
			Type:    formField.Type,
			Content: fieldContents[0],
		}

		buf, err := json.Marshal(subField.Content)
		if err != nil {
			return types.SubmissionData{}, fmt.Errorf("insert: failed to marshal field data: %v", err)

		}

		subField.ContentAsString = string(buf)

		for _, constraint := range formField.Constraints {
			if constraint.Name == "unique" {
				h := sha256.New()
				h.Write(buf)
				fieldHash := string(h.Sum(nil))

				exists, err := s.models.Submissions.CheckForRepeatedUniqueField(formInstanceID, fieldName, fieldHash)
				if err != nil {
					return types.SubmissionData{}, err
				}
				if exists {
					log.Printf("Insert: duplicate unique field detected: %v", fieldName)
					return types.SubmissionData{}, fmt.Errorf("models: duplicate unique field %v", fieldName)
				}

				subField.Unique = true
				subField.Hash = fieldHash
			}
		}

		submission.Fields = append(submission.Fields, subField)

		if len(fieldContents) > 1 {
			log.Printf("Insert: multiple values for field %s; only first saved", fieldName)
		}
	}

	return submission, nil
}
