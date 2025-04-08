package models

import (
	"database/sql"
	"errors"

	"formy.fprzg.net/internal/types"
	"formy.fprzg.net/internal/utils"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
	ErrInvalidInput       = errors.New("models: invalid input")
	ErrInvalidUserID      = errors.New("models: user not found")
	ErrUserNotFound       = errors.New("models: user not found")
	ErrFormNotFound       = errors.New("models: form not found")
)

const (
	ValidUserName     = "Alice"
	ValidUserEmail    = "alice@example.com"
	ValidUserPassword = "securepass"
)

type Models struct {
	Users       UsersModelInterface
	Forms       FormsModelInterface
	Submissions SubmissionsModelInterface
}

func GetModels(db *sql.DB) *Models {
	return &Models{
		Users:       &UsersModel{db},
		Forms:       &FormsModel{db},
		Submissions: &SubmissionsModel{db},
	}
}

// maybe we want to know the IDs of the user, form, form instance, submissions, etc
func GetTestModels() (*Models, error) {
	db, err := utils.SetupTestDB()
	if err != nil {
		panic(err)
	}

	m := GetModels(db)

	userID, err := InsertTestUser(m)
	if err != nil {
		return nil, err
	}

	_, err = InsertTestForms(m, userID)
	//formID, err := InsertTestForm(m, userID)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func InsertTestUser(m *Models) (int, error) {
	userID, err := m.Users.Insert(ValidUserName, ValidUserEmail, ValidUserPassword)
	if err != nil {
		return 0, err
	}

	_, err = m.Users.Authenticate(ValidUserEmail, ValidUserPassword)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func InsertTestForms(m *Models, userID int) ([]int, error) {
	form1Fields := []types.FieldData{
		{
			Name:        "name",
			Type:        "string",
			Constraints: []types.FieldConstraint{{Name: "required"}},
		},
		{
			Name:        "email",
			Type:        "string",
			Constraints: []types.FieldConstraint{{Name: "unique"}},
		},
	}

	form2Fields := []types.FieldData{
		{
			Name:        "name",
			Type:        "string",
			Constraints: []types.FieldConstraint{{Name: "required"}},
		},
		{
			Name:        "email",
			Type:        "string",
			Constraints: []types.FieldConstraint{{Name: "unique"}},
		},
	}

	formIDs := []int{0, 0}

	formID1, err := m.Forms.Insert(userID, "form1", "Form One", form1Fields)
	if err != nil {
		return nil, err
	}

	formID2, err := m.Forms.Insert(userID, "form2", "Form Two", form2Fields)
	if err != nil {
		return nil, err
	}

	formIDs[0] = formID1
	formIDs[0] = formID2

	return formIDs, nil
}
