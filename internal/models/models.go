package models

import (
	"database/sql"
	"errors"
	"time"

	"formy.fprzg.net/internal/types"
	"github.com/labstack/echo/v4"
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
	ValidUserName     = "alice"
	ValidUserPassword = "securepass"
)

type Models struct {
	Users           UsersModelInterface
	Forms           FormsModelInterface
	Submissions     SubmissionsModelInterface
	contextDuration time.Duration
}

var contextDuration time.Duration

func Get(db *sql.DB, e *echo.Echo, ctxDuration time.Duration) (*Models, error) {
	contextDuration = ctxDuration

	m := &Models{
		Users: &UsersModel{
			db: db,
			e:  e,
		},
		Forms: &FormsModel{
			db: db,
			e:  e,
		},
		Submissions: &SubmissionsModel{
			db: db,
			e:  e,
		},
	}

	return m, nil
}

func InsertTestUser(m *Models) (int, error) {
	userID, err := m.Users.Insert(ValidUserName, ValidUserPassword)
	if err != nil {
		return 0, err
	}

	_, err = m.Users.Authenticate(ValidUserName, ValidUserPassword)
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
