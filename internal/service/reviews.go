package service

import (
	"log/slog"

	"github.com/vermacodes/email-alert-processor/internal/entity"
)

func highVolumeReview(alert entity.Alert) ([]entity.AlertResponse, error) {
	slog.Info("Processing CaseHygiene alert")

	alertResponses := []entity.AlertResponse{}

	decodedAlertMessage, err := decodeBase64(alert.AlertMessage)
	if err != nil {
		slog.Error("Error decoding base64 string: ", err)
		return alertResponses, err
	}

	incidents, err := parseCaseBuddyEmail(decodedAlertMessage)
	if err != nil {
		slog.Error("Error parsing CaseBuddy alert: ", err)
		return alertResponses, err
	}

	slog.Info("Total CaseBuddy Reports: ", slog.Int("totalReports", len(incidents)))

	slog.Info("CaseBuddy Reports: ")
	slog.Info("------------------------------------------------")

	for _, incident := range incidents {
		slog.Info("IncidentOwnerAlias: ", slog.String("ownerAlias", incident.IncidentOwnerAlias))
		slog.Info("IncidentNumber: ", slog.String("caseNumber", incident.IncidentNumber))
		slog.Info(
			"IncidentCreatedTime: ",
			slog.String("caseCreatedTime", incident.IncidentCreatedTime),
		)
		slog.Info("IncidentState: ", slog.String("incidentState", incident.IncidentState))
		slog.Info("IncidentIdleDuration: ", slog.String("caseIdle", incident.IncidentIdleDuration))
		slog.Info("IncidentStatus: ", slog.String("caseStatus", incident.IncidentStatus))
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

func highImpactReview(alert entity.Alert) ([]entity.AlertResponse, error) {
	slog.Info("Processing CaseHygiene alert")

	alertResponses := []entity.AlertResponse{}

	decodedAlertMessage, err := decodeBase64(alert.AlertMessage)
	if err != nil {
		slog.Error("Error decoding base64 string: ", err)
		return alertResponses, err
	}

	incidents, err := parseCaseBuddyEmail(decodedAlertMessage)
	if err != nil {
		slog.Error("Error parsing CaseBuddy alert: ", err)
		return alertResponses, err
	}

	slog.Info("Total incidents in alert message:", slog.Int("incidentCount", len(incidents)))

	slog.Info("Incidnet Details")
	slog.Info("------------------------------------------------")

	for _, incident := range incidents {
		slog.Info("IncidentOwnerAlias: ", slog.String("ownerAlias", incident.IncidentOwnerAlias))
		slog.Info("IncidentNumber: ", slog.String("caseNumber", incident.IncidentNumber))
		slog.Info(
			"IncidentCreatedTime: ",
			slog.String("caseCreatedTime", incident.IncidentCreatedTime),
		)
		slog.Info("IncidentState: ", slog.String("incidentState", incident.IncidentState))
		slog.Info("IncidentIdleDuration: ", slog.String("caseIdle", incident.IncidentIdleDuration))
		slog.Info("IncidentStatus: ", slog.String("caseStatus", incident.IncidentStatus))
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
