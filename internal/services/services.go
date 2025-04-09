package services

import "formy.fprzg.net/internal/models"

type Services struct {
	Submission SubmissionServiceInterface
}

func Get(m *models.Models) *Services {
	return &Services{
		Submission: &SubmissionService{m},
	}
}
