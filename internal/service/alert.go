package service

import (
	"github.com/vermacodes/email-alert-processor/internal/entity"
)

type alertService struct{}

func NewAlertService() entity.AlertService {
	return &alertService{}
}

func (s *alertService) ProcessAlert(alert entity.Alert) (entity.Alert, error) {
	switch alert.ID {
	case "Confluent":
		return cnflntNewCaseAlert(alert)
	case "CaseHygiene":
		return caseHygiene(alert)
	case "HighImpact":
		return highImpactReview(alert)
	case "HighVolume":
		return highVolumeReview(alert)
	default:
		return alert, nil
	}
}

func addTwoNumbers(a, b int) int {
	return a + b
}
