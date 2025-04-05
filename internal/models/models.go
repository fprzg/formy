package models

import (
	"database/sql"
	"errors"
	"testing"

	"formy.fprzg.net/internal/utils"
	"github.com/stretchr/testify/assert"
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

func GetModels(db *sql.DB) Models {
	return Models{
		Users:       &UsersModel{db},
		Forms:       &FormsModel{db},
		Submissions: &SubmissionsModel{db},
	}
}

// Should this one insert a lil bit of data (a user, two forms, some submitts)
func setupTestDB(t *testing.T) Models {
	ctx, err := utils.NewMigrationCtx(":memory:", ":memory:", "../../migrations")
	assert.NoError(t, err)

	err = ctx.Migrate(-1)
	assert.NoError(t, err)

	m := GetModels(ctx.AppDB)
	_ = insertTestUser(t, m, ValidUserName, ValidUserEmail, ValidUserPassword)

	const formFields = `[ {"field_name": "email", "field_type": "string", "contraints": ["unique"]} ]`
	err = m.Forms.InsertForm(1, "form name", "form description", formFields)
	assert.NoError(t, err)

	ctx.StateDB.Close()

	return m
}

func insertTestUser(t *testing.T, m Models, name, email, password string) int {
	err := m.Users.Insert(name, email, password)
	assert.NoError(t, err)

	id, err := m.Users.Authenticate(email, password)
	assert.NoError(t, err)

	return id
}
