package models

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
	LastUpdated time.Time `json:"last_updated"`
	LastLogin   time.Time `json:"last_login"`
}

type userData struct {
	id           int
	passwordHash []byte
	*User
}

type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
	Get(id int) (User, error)
	UpdateName(id int, name, password string) error
	UpdateEmail(id int, email, password string) error
	UpdatePassword(id int, oldPwd, newPwd string) error
}

type UserModel struct {
	db *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	const query = `
	INSERT INTO users (name, email, password)
	VALUES (?, ?, ?)
	RETURNING id, name, email, created_at, last_updated`

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	var u User
	var id int
	err = m.db.QueryRow(query, name, email, passwordHash).Scan(&id, &u.Name, &u.Email, &u.CreatedAt, &u.LastUpdated)
	return err
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	const query = `
	SELECT id, password FROM users WHERE email = ?`

	var u userData
	err := m.db.QueryRow(query, email).Scan(&u.id, &u.passwordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(u.passwordHash, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	return u.id, nil
}

func (m *UserModel) AuthenticateUsingID(id int, password string) error {
	const query = `
	SELECT password FROM users WHERE id = ?
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

func (m *UserModel) Exists(id int) (bool, error) {
	const query = `
	SELECT EXISTS(SELECT true FROM users WHERE id = ?)
	`

	var exists bool
	err := m.db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, ErrInvalidCredentials
		}
		return false, nil
	}

	return exists, nil
}

func (m *UserModel) Get(id int) (User, error) {
	query := `
	SELECT name, email, created_at, last_updated, last_login
	FROM users
	WHERE id = ?
	`

	var u User
	err := m.db.QueryRow(query, id).Scan(&u.Name, &u.Email, &u.CreatedAt, &u.LastUpdated, &u.LastLogin)
	if err == sql.ErrNoRows {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrInvalidCredentials
		}
		return User{}, nil
	}

	return u, err
}

func (m *UserModel) UpdateName(id int, name, password string) error {
	const query = `
	UPDATE users SET name = ? WHERE id = ?
	`

	err := m.AuthenticateUsingID(id, password)
	if err != nil {
		return err
	}

	_, err = m.db.Exec(query, name, id)

	return err
}

func (m *UserModel) UpdateEmail(id int, email, password string) error {
	const query = `
	UPDATE users SET email = ? WHERE id = ?
	`

	err := m.AuthenticateUsingID(id, password)
	if err != nil {
		return err
	}

	_, err = m.db.Exec(query, email, id)

	return err
}

func (m *UserModel) UpdatePassword(id int, oldPwd, newPwdRaw string) error {
	const query = `
	UPDATE users SET password = ? WHERE id = ?
	`

	err := m.AuthenticateUsingID(id, oldPwd)
	if err != nil {
		return err
	}

	newPwd, err := bcrypt.GenerateFromPassword([]byte(newPwdRaw), 12)
	if err != nil {
		return err
	}

	_, err = m.db.Exec(query, string(newPwd), id)

	return err
}
