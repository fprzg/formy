package services

import (
	"formy.fprzg.net/internal/models"
	"github.com/labstack/echo/v4"
)

type Services struct {
	jwtSecret       string
	models          *models.Models
	e               *echo.Echo
	TemplateManager *TemplateManager
}

func Get(jwtSecret string, m *models.Models, tm *TemplateManager, e *echo.Echo) (*Services, error) {
	return &Services{
		jwtSecret:       jwtSecret,
		models:          m,
		e:               e,
		TemplateManager: tm,
	}, nil
}

func (s *Services) Close() {
	s.TemplateManager.Close()
}
