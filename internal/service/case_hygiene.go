package service

import (
	"log/slog"
	"strings"

	"golang.org/x/net/html"

	"github.com/vermacodes/email-alert-processor/internal/entity"
)

func caseHygiene(alert entity.Alert) ([]entity.AlertResponse, error) {
	slog.Info("Processing CaseHygiene alert")

	alertResponses := []entity.AlertResponse{}

	decodedAlertMessage, err := decodeBase64(alert.AlertMessage)
	if err != nil {
		slog.Error("Error decoding base64 string: ", err)
		return alertResponses, err
	}

	incidents, err := parseCaseBuddyAlert(decodedAlertMessage)
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

func parseCaseBuddyAlert(alert string) ([]entity.Incident, error) {
	doc, err := html.Parse(strings.NewReader(alert))
	if err != nil {
		return []entity.Incident{}, err
	}

	var incidents []entity.Incident
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {
			var incident entity.Incident
			tdCount := 0
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "td" {
					var data string
					if c.FirstChild.Data == "a" {
						data = c.FirstChild.FirstChild.Data
					} else {
						data = c.FirstChild.Data
					}
					switch tdCount {
					case 1:
						incident.IncidentNumber = data
					case 2:
						incident.IncidentOwnerAlias = data
					case 3:
						incident.IncidentCreatedTime = data
					case 4:
						incident.IncidentState = data
					case 5:
						incident.IncidentIdleDuration = data
					case 6:
						incident.IncidentStatus = data
					}
					tdCount++

					slog.Info("Data: ", slog.String("data", data))
				}
			}
			if tdCount > 0 {
				incidents = append(incidents, incident)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return incidents, nil
}
