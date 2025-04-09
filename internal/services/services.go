package services

import "formy.fprzg.net/internal/models"

type Services struct {
	Submissions SubmissionsServiceInterface
	Forms       FormsServiceInterface
}

func Get(m *models.Models) *Services {
	return &Services{
		Submissions: &SubmissionsService{m},
		Forms:       &FormsService{m},
	}
}
