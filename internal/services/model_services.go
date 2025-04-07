package services

import (
	"errors"

	"formy.fprzg.net/internal/models"
)

type ModelServices struct {
	FormsServices *FormsServices
}

var (
	ErrInvalidInput = errors.New("service: Invalid input")
)

func GetModelServices(m *models.Models) *ModelServices {
	return &ModelServices{
		FormsServices: &FormsServices{fm: m.Forms, fim: m.FormInstances},
	}
}
