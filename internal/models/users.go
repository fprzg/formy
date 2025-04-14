package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"formy.fprzg.net/internal/utils"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserName    string    `json:"user_name"`
	CreatedAt   time.Time `json:"created_at"`
	LastUpdated time.Time `json:"last_updated"`
	LastLogin   time.Time `json:"last_login"`
}

type userData struct {
	ID           int
	PasswordHash []byte
	User
}

type UsersModelInterface interface {
	Insert(nameName, password string) (int, error)
	Authenticate(userName, password string) (int, error)
	Exists(id int) (bool, error)
	Get(id int) (User, error)
	UpdatePassword(id int, oldPwd, newPwd string) error
}

type UsersModel struct {
	db *sql.DB
	e  *echo.Echo
}

func (m *UsersModel) Insert(userName, password string) (int, error) {
	if userName == "" || password == "" {
		return 0, ErrInvalidInput
	}

	const query = `
	INSERT INTO users (user_name, password)
	VALUES (?, ?)
	RETURNING id, created_at, updated_at
	`

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, err
	}

	var u userData
	u.UserName = userName
	err = m.db.QueryRow(query, userName, passwordHash).Scan(&u.ID, &u.CreatedAt, &u.LastUpdated)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.user_name") {
			return 0, ErrDuplicateEmail
		}
		return 0, err
	}

	return u.ID, nil
}

func (m *UsersModel) Authenticate(userName, password string) (int, error) {
	const query = `
	SELECT id, password
	FROM users
	WHERE user_name = ?
	`

	var u userData
	err := m.db.QueryRow(query, userName).Scan(&u.ID, &u.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	return u.ID, nil
}

func (m *UsersModel) AuthenticateUsingID(id int, password string) error {
	const query = `
	SELECT password
	FROM users
	WHERE id = ?
	`

	var pwd []byte
	err := m.db.QueryRow(query, id).Scan(&pwd)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInvalidCredentials
		}
		return err
	}

	err = bcrypt.CompareHashAndPassword(pwd, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
		}
		return err
	}

	return nil
}

func (m *UsersModel) Exists(id int) (bool, error) {
	const query = `
	SELECT EXISTS(
		SELECT true
		FROM users
		WHERE id = ?
	)
	`

	var exists bool
	err := m.db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, ErrNoRecord
		}
		return false, nil
	}

	return exists, nil
}

func (m *UsersModel) Get(id int) (User, error) {
	const query = `
	SELECT user_name, created_at, updated_at, last_login
	FROM users
	WHERE id = ?
	`

	var u User
	err := m.db.QueryRow(query, id).Scan(&u.UserName, &u.CreatedAt, &u.LastUpdated, &u.LastLogin)
	if err == sql.ErrNoRows {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNoRecord
		}
		return User{}, nil
	}

	return u, err
}

func (m *UsersModel) UpdatePassword(id int, oldPwd, newPwdRaw string) error {
	if newPwdRaw == "" {
		return ErrInvalidInput
	}

	const query = `
	UPDATE users
	SET password = ?
	WHERE id = ?
	`

	err := m.AuthenticateUsingID(id, oldPwd)
	if err != nil {
		return err
	}

	newPwd, err := bcrypt.GenerateFromPassword([]byte(newPwdRaw), 12)
	if err != nil {
		return err
	}

	rows, err := utils.ExecuteSqlStmt(m.db, query, string(newPwd), id)
	if rows == 0 {
		return ErrUserNotFound
	}

	return err
}
