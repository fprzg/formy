package models

import "database/sql"

type FormInstancesModelInterface interface {
	Insert() error
}

type FormInstances struct {
	db *sql.DB
}

func (m *FormInstances) Insert() error {
	return nil
}
