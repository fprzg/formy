package models

import (
	"database/sql"
	"errors"

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
	Users         UsersModelInterface
	Forms         FormsModelInterface
	FormInstances FormInstancesModelInterface
	Submissions   SubmissionsModelInterface
}

func GetModels(db *sql.DB) *Models {
	return &Models{
		Users:         &UsersModel{db},
		Forms:         &FormsModel{db},
		FormInstances: &FormInstances{db},
		Submissions:   &SubmissionsModel{db},
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

	_, err = InsertTestForm(m, userID)
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

func InsertTestForm(m *Models, userID int) (int, error) {
	formID, err := m.Forms.Insert(userID, "form name", "form description", `[ {"field_name": "email", "field_type": "string", "field_contraints": ["unique"]} ]`)
	return formID, err
}
