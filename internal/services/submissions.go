package services

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"formy.fprzg.net/internal/models"
	"formy.fprzg.net/internal/types"
	"github.com/labstack/echo/v4"
)

type SubmissionsService struct {
	models *models.Models
	e      *echo.Echo
}

type SubmissionsServiceInterface interface {
	ProcessSubmission(formID int, r *http.Request, ctx context.Context) (int, error)
	ParseSubmissionFromRequest(form types.FormData, r *http.Request, ctx context.Context) (types.SubmissionData, error)
}

func (s *SubmissionsService) ProcessSubmission(formID int, r *http.Request, ctx context.Context) (int, error) {
	formData, err := s.models.Forms.Get(formID)
	if err != nil {
		return 0, err
	}

	submission, err := s.ParseSubmissionFromRequest(formData, r, ctx)
	if err != nil {
		return 0, err
	}

	return s.models.Submissions.Insert(submission, ctx)
}

func (s *SubmissionsService) ParseSubmissionFromRequest(form types.FormData, r *http.Request, ctx context.Context) (types.SubmissionData, error) {
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
			s.e.Logger.Printf("Insert: unknown field %s; skipping.\n", fieldName)
			continue
		}

		formField := form.Fields[fieldIndex]

		subField := types.SubmissionField{
			Name:    formField.Name,
			Type:    formField.Type,
			Content: fieldContents[0],
		}

		if !types.TypesMatch(subField.Content, subField.Type) {
			receivedType := reflect.ValueOf(subField.Content).Kind()
			return types.SubmissionData{}, fmt.Errorf("invalid field type at '%s': expected '%s' but received '%v'", subField.Name, subField.Type, receivedType)
		}

		if subField.Type == "string" {
			subField.ContentAsString = subField.Content.(string)
		} else {

			buf, err := json.Marshal(subField.Content)
			if err != nil {
				return types.SubmissionData{}, fmt.Errorf("insert: failed to marshal field data: %v", err)
			}
			subField.ContentAsString = string(buf)
		}

		for _, constraint := range formField.Constraints {
			if constraint.Name == "unique" {
				h := sha256.New()
				h.Write([]byte(subField.ContentAsString))
				fieldHash := string(h.Sum(nil))

				exists, err := s.models.Submissions.CheckForRepeatedUniqueField(formInstanceID, fieldName, fieldHash)
				if err != nil {
					return types.SubmissionData{}, err
				}
				if exists {
					s.e.Logger.Printf("Insert: duplicate unique field detected: '%s'.\n", fieldName)
					return types.SubmissionData{}, fmt.Errorf("models: duplicate unique field %v", fieldName)
				}

				subField.Unique = true
				subField.Hash = fieldHash
			}
		}

		submission.Fields = append(submission.Fields, subField)

		if len(fieldContents) > 1 {
			s.e.Logger.Printf("Insert: multiple values for field %s; only first saved.\n", fieldName)
		}
	}

	return submission, nil
}
