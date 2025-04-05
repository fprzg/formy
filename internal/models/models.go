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
func setupTestDB(t *testing.T) (utils.MigrationCtx, Models) {
	ctx, err := utils.NewMigrationCtx(":memory:", ":memory:", "../../migrations")
	assert.NoError(t, err)

	err = ctx.Migrate(-1)
	assert.NoError(t, err)

	m := GetModels(ctx.AppDB)
	_ = insertTestUser(t, m, ValidUserName, ValidUserEmail, ValidUserPassword)

	// ?????
	//ctx.StateDB.Close()

	return ctx, m
}

func insertTestUser(t *testing.T, m Models, name, email, password string) int {
	err := m.Users.Insert(name, email, password)
	assert.NoError(t, err)

	id, err := m.Users.Authenticate(email, password)
	assert.NoError(t, err)

	return id
}
