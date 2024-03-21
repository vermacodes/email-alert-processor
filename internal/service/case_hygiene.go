package service

import (
	"encoding/base64"
	"log/slog"

	"github.com/vermacodes/email-alert-processor/internal/entity"
)

func caseHygiene(alert entity.Alert) ([]entity.AlertResponse, error) {
	slog.Info("Processing CaseHygiene alert")

	alertResponses := []entity.AlertResponse{}

	decodedAlertMessage, err := base64.StdEncoding.DecodeString(alert.AlertMessage)
	if err != nil {
		slog.Error("Error decoding base64 string: ", err)
		return alertResponses, err
	}

	incidents, err := parseCaseBuddyAlert(string(decodedAlertMessage))
	if err != nil {
		slog.Error("Error parsing CaseBuddy alert: ", err)
		return alertResponses, err
	}

	slog.Info("Total CaseBuddy incidents: ", slog.Int("totalincidents", len(incidents)))

	slog.Info("CaseBuddy incidents: ")
	slog.Info("------------------------------------------------")

	for _, incident := range incidents {
		slog.Info(
			"IncidentOwnerAlias: ",
			slog.String("IncidentOwnerAlias", incident.IncidentOwnerAlias),
		)
		slog.Info("IncidentNumber: ", slog.String("IncidentNumber", incident.IncidentNumber))
		slog.Info(
			"IncidentCreatedTime: ",
			slog.String("IncidentCreatedTime", incident.IncidentCreatedTime),
		)
		slog.Info("IncidentState: ", slog.String("IncidentState", incident.IncidentState))
		slog.Info(
			"IncidentIdleDuration: ",
			slog.String("IncidentIdleDuration", incident.IncidentIdleDuration),
		)
		slog.Info("IncidentStatus: ", slog.String("IncidentStatus", incident.IncidentStatus))
		slog.Info("------------------------------------------------")

		alertResponse := entity.AlertResponse{
			AlertID:      alert.ID,
			AlertType:    alert.AlertType,
			AlertMessage: alert.AlertMessage,
			Incident:     incident,
		}

		alertResponses = append(alertResponses, alertResponse)
	}

	return alertResponses, nil
}
