package services

import (
	"formy.fprzg.net/internal/models"
	"github.com/labstack/echo/v4"
)

type Services struct {
	Submissions     SubmissionsServiceInterface
	Forms           FormsServiceInterface
	TemplateManager *TemplateManager
}

func Get(m *models.Models, tm *TemplateManager, e *echo.Echo) (*Services, error) {
	return &Services{
		Submissions:     &SubmissionsService{m, e},
		Forms:           &FormsService{m, e},
		TemplateManager: tm,
	}, nil
}

func (s *Services) Close() {
	s.TemplateManager.Close()
}
