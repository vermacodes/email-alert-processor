package service

import (
	"encoding/base64"
	"log/slog"

	"github.com/vermacodes/email-alert-processor/internal/entity"
)

func caseBuddyAlert(alert entity.Alert) ([]entity.AlertResponse, error) {
	decodedEmail, err := base64.StdEncoding.DecodeString(alert.AlertMessage)
	if err != nil {
		slog.Error("Error decoding email body",
			slog.String("alertId", alert.ID),
			slog.String("error", err.Error()),
		)
		return []entity.AlertResponse{}, err
	}

	incidents, err := parseCaseBuddyEmail(string(decodedEmail))
	if err != nil {
		slog.Error("Error parsing case buddy alert",
			slog.String("alertId", alert.ID),
			slog.String("error", err.Error()),
		)
		return []entity.AlertResponse{}, err
	}

	var alertResponses []entity.AlertResponse
	for _, incident := range incidents {
		alertResponses = append(alertResponses, entity.AlertResponse{
			AlertID:  alert.ID,
			Incident: incident,
		})

		slog.Info("Incident: ", slog.String("incident", incident.IncidentNumber))
	}

	return alertResponses, nil
}
