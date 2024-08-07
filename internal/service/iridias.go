package service

import (
	"encoding/base64"
	"log/slog"

	"github.com/vermacodes/email-alert-processor/internal/entity"
)

func iridiasNewCaseAlert(alert entity.Alert) ([]entity.AlertResponse, error) {
	// Process alert here.
	// Alert message is base64 encoded HTML content. Decode it first.

	slog.Info("Iridias New Case Alert")

	decodedAlertMessage, err := base64.StdEncoding.DecodeString(alert.AlertMessage)
	if err != nil {
		slog.Error("Error decoding base64 string: ", err)
		return []entity.AlertResponse{}, err
	}

	extractedText, err := extractTextFromIridiasEmail(string(decodedAlertMessage))
	if err != nil {
		slog.Error("Error extracting text from Iridias email: ", err)
		return []entity.AlertResponse{}, err
	}

	// Extract the incident from the text.
	incident, err := parseTextFromIridiasEmail(extractedText)
	if err != nil {
		slog.Error("Error extracting incident from text: ", err)
		return []entity.AlertResponse{}, err
	}

	alertResponse := entity.AlertResponse{
		AlertID:      alert.ID,
		AlertType:    alert.AlertType,
		AlertMessage: "",
		Incident:     incident,
	}

	alertResponses := []entity.AlertResponse{alertResponse}

	return alertResponses, nil
}
