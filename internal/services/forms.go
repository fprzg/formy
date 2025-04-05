package services

import (
	"encoding/json"

	"formy.fprzg.net/internal/models"
	"formy.fprzg.net/internal/types"
)

type FormServiceInterface interface {
	CreateForm(userID int, form types.JSONMap) error
	UpdateForm(formID int, formUpdate types.JSONMap) error
	SubmitForm(formID int, fieldValues types.JSONMap) error
	//GetSubmissions(formID int) ([]FormInstance, error)
}

type FormService struct {
	m models.FormsModelInterface
}

func (s *FormService) CreateForm(userID int, form types.JSONMap) error {
	name, ok := form["name"].(string)
	if !ok {
		return nil
	}
	description, ok := form["description"].(string)
	if !ok {
		return nil
	}
	fieldsMap, ok := form["fields"]
	if !ok {
		return nil
	}

	bytes, err := json.Marshal(fieldsMap)
	if err != nil {
		return err
	}

	fields := string(bytes)

	err = s.m.InsertForm(userID, name, description, fields)

	return err
}

func (s *FormService) UpdateForm(formID int, formUpdates types.JSONMap) error {
	return nil
}

func (s *FormService) SubmitForm(formID int, fieldValues types.JSONMap) error {
	return nil
}

/*
func (s *FormService) GetSubmissions(formID int) ([]FormInstance, error) {
	return nil, nil
}
*/
