package models

import (
	"database/sql"
	"errors"
	"log"
	"time"

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

func GetModels(cfg types.AppConfig) (*Models, error) {
	db, duration, err := DatabaseSetup(cfg)
	return models(db, duration), err
}

func DatabaseSetup(cfg types.AppConfig) (*sql.DB, time.Duration, error) {
	var db *sql.DB
	var duration time.Duration
	var err error

	if cfg.Env == "development" || cfg.Env == "testing" {
		duration = 1 * time.Hour

		db, err = utils.SetupTestDB()
		if err != nil {
			log.Fatal(err.Error())
		}

		m := models(db, duration)

		userID, err := InsertTestUser(m)
		if err != nil {
			log.Fatal(err.Error())
		}

		_, err = InsertTestForms(m, userID)
		if err != nil {
			log.Fatal(err.Error())
		}
	} else {
		duration = 3 * time.Second
		db, err = sql.Open("sqlite3", cfg.DBDir)
	}

	return db, duration, err
}

func models(db *sql.DB, duration time.Duration) *Models {
	return &Models{
		Users: &UsersModel{
			db:                  db,
			transactionDuration: duration,
		},
		Forms: &FormsModel{
			db:                  db,
			transactionDuration: duration,
		},
		Submissions: &SubmissionsModel{
			db:                  db,
			transactionDuration: duration,
		},
	}
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
	form1Fields := []types.FormField{
		{
			Name:        "name",
			Type:        "string",
			Constraints: []types.FieldConstraint{{Name: "required"}},
		},
		{
			Name:        "email",
			Type:        "string",
			Constraints: []types.FieldConstraint{{Name: "email"}, {Name: "unique"}},
		},
		{
			Name:        "subject",
			Type:        "string",
			Constraints: nil,
		},
		{
			Name:        "message",
			Type:        "string",
			Constraints: nil,
		},
	}

	form2Fields := []types.FormField{
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

	fid1, err := m.Forms.Insert(userID, "hapaxredux.com contact form", "Contact form for hapaxredux.com CTA.", form1Fields)
	if err != nil {
		return nil, err
	}

	fid2, err := m.Forms.Insert(userID, "form2", "Form Two", form2Fields)
	if err != nil {
		return nil, err
	}

	return []int{fid1, fid2}, nil
}
