package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"formy.fprzg.net/internal/models"
	"formy.fprzg.net/internal/types"
	"github.com/labstack/echo/v4"
)

type FormsService struct {
	models *models.Models
	e      *echo.Echo
}

type FormsServiceInterface interface {
	ProcessForm(r *http.Request, ctx context.Context) (int, error)
	ParseFormFromRequest(r *http.Request, ctx context.Context) (types.FormData, error)
}

func (s *FormsService) ProcessForm(r *http.Request, ctx context.Context) (int, error) {
	formData, err := s.ParseFormFromRequest(r, r.Context())
	if err != nil {
		return 0, err
	}

	formID, err := s.models.Forms.Insert(formData.UserID, formData.Name, formData.Description, formData.Fields)
	if err != nil {
		return 0, err
	}

	return formID, nil
}

func (s *FormsService) ParseFormFromRequest(r *http.Request, ctx context.Context) (types.FormData, error) {
	if err := r.ParseForm(); err != nil {
		return types.FormData{}, err
	}

	userID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		return types.FormData{}, err
	}

	formData := types.FormData{
		UserID:      userID,
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
	}

	fieldNames := r.Form["field_name"]
	fieldTypes := r.Form["field_type"]
	fieldConstraintsString := r.Form["field_constraints"]

	if len(fieldNames) == 0 || len(fieldNames) != len(fieldTypes) || len(fieldNames) != len(fieldConstraintsString) {
		return types.FormData{}, fmt.Errorf("invalid fields data")
	}

	for i := range fieldNames {
		var fieldConstraints []types.FieldConstraint
		err = json.Unmarshal([]byte(fieldConstraintsString[i]), &fieldConstraints)
		if err != nil {
			return types.FormData{}, err
		}

		formData.Fields = append(formData.Fields, types.FormField{
			Name:        fieldNames[i],
			Type:        fieldTypes[i],
			Constraints: fieldConstraints,
		})
	}

	return formData, nil
}
