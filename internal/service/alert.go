package service

import (
	"log/slog"

	"github.com/vermacodes/email-alert-processor/internal/entity"
)

type alertService struct{}

func NewAlertService() entity.AlertService {
	return &alertService{}
}

func (s *alertService) ProcessAlert(alert entity.Alert) ([]entity.AlertResponse, error) {
	switch alert.ID {
	case "Iridias":
		return iridiasNewCaseAlert(alert)
	case "CaseHygiene":
		return caseHygiene(alert)
	case "HighImpact":
		return highImpactReview(alert)
	case "HighVolume":
		return highVolumeReview(alert)
	default:
		slog.Error("Not a valid alert type",
			slog.String("alertId", alert.ID),
		)
		return []entity.AlertResponse{}, nil
	}
}

func addTwoNumbers(a, b int) int {
	return a + b
}
