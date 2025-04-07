package services

import (
	"formy.fprzg.net/internal/models"
	"formy.fprzg.net/internal/types"
)

type FormsServicesInterface interface {
	CreateForm(userID int, form types.JSONMap) error
	UpdateForm(formID int, formUpdate types.JSONMap) error
	SubmitForm(formID int, fieldValues types.JSONMap) error
	//GetSubmissions(formID int) ([]FormInstance, error)
}

type FormsServices struct {
	fm models.FormsModelInterface
}

/*
func (s *FormsServices) CreateForm(form types.FormData) (int, error) {
	// create the form
	jsonByte, err := json.Marshal(form.Fields)
	if err != nil {
		return 0, err
	}

	//json := string(jsonByte)
	userID, err := strconv.Atoi(form.UserID)
	if err != nil {
		return 0, err
	}

	formID, err := s.fm.Insert(userID, form.Name, form.Description, json)
	if err != nil {
		return 0, err
	}

	instanceID, err := s.fim.Insert(formID, form.Fields)
	utils.Print(instanceID)

	return formID, nil
}

func (s *FormsServices) CreateForm(userID int, form types.JSONMap) error {
	var err error
	fields := ""

	name, ok := form["name"].(string)
	if !ok {
		return ErrInvalidInput
	}
	description, ok := form["description"].(string)
	if !ok {
		return ErrInvalidInput
	}

	fieldsArray, ok := form["fields"].([]interface{})
	if !ok {
		return ErrInvalidInput
	}

	for _, fieldAsInterface := range fieldsArray {
		field, ok := fieldAsInterface.(map[string]interface{})
		if !ok {
			return ErrInvalidInput
		}

		fieldName, ok := field["field_name"].(string)
		if !ok {
			return ErrInvalidInput
		}
		fieldType, ok := field["field_type"].(string)
		if !ok {
			return ErrInvalidInput
		}
		fieldConstraints, ok := field["field_constraints"]
		if !ok {
			// no problem hehe
		}
		utils.Print(fieldConstraints)

		err := ValidateName(fieldName)
		if err != nil {
			return err
		}

		err = ValidateType(fieldType)
		if err != nil {
			return err
		}

		err = ValidateConstraints(fieldType, fieldConstraints)
		if err != nil {
			return err
		}
	}

	err = s.fm.Insert(userID, name, description, fields)

	return err
}
*/

func (s *FormsServices) UpdateForm(formID int, formUpdates types.FormData) error {
	return nil
}

func (s *FormsServices) SubmitForm(formID int, fieldValues types.FormData) error {
	return nil
}

/*
func (s *FormsServices) GetSubmissions(formID int) ([]FormInstance, error) {
	return nil, nil
}
*/
