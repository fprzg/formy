package models

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	schema := `
	CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_updated DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_login DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT 0,
    token_version INTEGER DEFAULT 0
	);`

	_, err = db.Exec(schema)
	assert.NoError(t, err)

	return db
}

func insertTestUser(t *testing.T, model *UserModel, name, email, password string) int {
	err := model.Insert(name, email, password)
	assert.NoError(t, err)

	id, err := model.Authenticate(email, password)
	assert.NoError(t, err)

	return id
}

func TestInsertAndAuthenticate(t *testing.T) {
	db := setupTestDB(t)
	model := &UserModel{db: db}

	err := model.Insert("Alice", "alice@example.com", "securepass")
	assert.NoError(t, err)

	id, err := model.Authenticate("alice@example.com", "securepass")
	assert.NoError(t, err)
	assert.NotZero(t, id)

	_, err = model.Authenticate("alice@example.com", "wrongpass")
	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestExists(t *testing.T) {
	db := setupTestDB(t)
	model := &UserModel{db: db}

	id := insertTestUser(t, model, "Bob", "bob@example.com", "pass123")

	exists, err := model.Exists(id)
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = model.Exists(999)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestGet(t *testing.T) {
	db := setupTestDB(t)
	model := &UserModel{db: db}

	id := insertTestUser(t, model, "Carol", "carol@example.com", "secret")

	user, err := model.Get(id)
	assert.NoError(t, err)
	assert.Equal(t, "Carol", user.Name)
	assert.Equal(t, "carol@example.com", user.Email)
}

func TestUpdateName(t *testing.T) {
	db := setupTestDB(t)
	model := &UserModel{db: db}

	id := insertTestUser(t, model, "Dave", "dave@example.com", "pass")

	err := model.UpdateName(id, "Dave Updated", "pass")
	assert.NoError(t, err)

	user, err := model.Get(id)
	assert.NoError(t, err)
	assert.Equal(t, "Dave Updated", user.Name)
}

func TestUpdateEmail(t *testing.T) {
	db := setupTestDB(t)
	model := &UserModel{db: db}

	id := insertTestUser(t, model, "Eve", "eve@old.com", "pass")

	err := model.UpdateEmail(id, "eve@new.com", "pass")
	assert.NoError(t, err)

	_, err = model.Authenticate("eve@new.com", "pass")
	assert.NoError(t, err)
}

func TestUpdatePassword(t *testing.T) {
	db := setupTestDB(t)
	model := &UserModel{db: db}

	id := insertTestUser(t, model, "Frank", "frank@example.com", "oldpass")

	err := model.UpdatePassword(id, "oldpass", "newpass")
	assert.NoError(t, err)

	_, err = model.Authenticate("frank@example.com", "oldpass")
	assert.ErrorIs(t, err, ErrInvalidCredentials)

	_, err = model.Authenticate("frank@example.com", "newpass")
	assert.NoError(t, err)
}
