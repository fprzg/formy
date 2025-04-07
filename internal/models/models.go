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
	Users         UsersModelInterface
	Forms         FormsModelInterface
	FormInstances FormInstancesModelInterface
	Submissions   SubmissionsModelInterface
}

func ExecuteSqlStmt(db *sql.DB, stmt string, args ...any) (int64, error) {
	result, err := db.Exec(stmt, args...)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func GetModels(db *sql.DB) Models {
	return Models{
		Users:         &UsersModel{db},
		Forms:         &FormsModel{db},
		FormInstances: &FormInstances{db},
		Submissions:   &SubmissionsModel{db},
	}
}

func SetupTestDB(t *testing.T) Models {
	ctx, err := utils.NewMigrationCtx(":memory:", ":memory:", "../../migrations")
	assert.NoError(t, err)

	err = ctx.Migrate(-1)
	assert.NoError(t, err)

	m := GetModels(ctx.AppDB)
	_ = InsertTestUser(t, m, ValidUserName, ValidUserEmail, ValidUserPassword)

	const formFields = `[ {"field_name": "email", "field_type": "string", "contraints": ["unique"]} ]`
	err = m.Forms.Insert(1, "form name", "form description", formFields)
	assert.NoError(t, err)

	ctx.StateDB.Close()

	return m
}

func InsertTestUser(t *testing.T, m Models, name, email, password string) int {
	err := m.Users.Insert(name, email, password)
	assert.NoError(t, err)

	id, err := m.Users.Authenticate(email, password)
	assert.NoError(t, err)

	return id
}
