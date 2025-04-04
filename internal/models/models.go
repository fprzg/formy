package models

import (
	"database/sql"
	"errors"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

type Models struct {
	formsModel       FormsModelInterface
	submissionsModel SubmissionsModelInterface
}

func GetModels(db *sql.DB) Models {
	return Models{
		formsModel:       &FormsModel{db},
		submissionsModel: &SubmissionsModel{db},
	}
}
