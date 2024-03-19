package service

import (
	"log/slog"
	"strings"

	cards "github.com/DanielTitkov/go-adaptive-cards"
	"golang.org/x/net/html"

	"github.com/vermacodes/email-alert-processor/internal/entity"
)

func caseHygiene(alert entity.Alert) (entity.Alert, error) {
	slog.Info("Processing CaseHygiene alert")

	decodedAlertMessage, err := decodeBase64(alert.AlertMessage)
	if err != nil {
		slog.Error("Error decoding base64 string: ", err)
		return alert, err
	}

	caseBuddyReports, err := parseCaseBuddyAlert(decodedAlertMessage)
	if err != nil {
		slog.Error("Error parsing CaseBuddy alert: ", err)
		return alert, err
	}

	slog.Info("Total CaseBuddy Reports: ", slog.Int("totalReports", len(caseBuddyReports)))

	slog.Info("CaseBuddy Reports: ")
	slog.Info("------------------------------------------------")

	for _, report := range caseBuddyReports {
		slog.Info("OwnerAlias: ", slog.String("ownerAlias", report.OwnerAlias))
		slog.Info("CaseNumber: ", slog.String("caseNumber", report.CaseNumber))
		slog.Info("CaseCreatedTime: ", slog.String("caseCreatedTime", report.CaseCreatedTime))
		slog.Info("IncidentState: ", slog.String("incidentState", report.IncidentState))
		slog.Info("CaseIdle: ", slog.String("caseIdle", report.CaseIdle))
		slog.Info("CaseStatus: ", slog.String("caseStatus", report.CaseStatus))
		slog.Info("------------------------------------------------")

		adaptiveCard := buildAdaptiveCard(report)
		alert.TeamsMessages = append(alert.TeamsMessages, entity.TeamsMessage{
			AudType:      "user",
			AudId:        "ashisverma",
			AdaptiveCard: adaptiveCard,
		})
	}

	return alert, nil
}

func buildAdaptiveCard(caseBuddyReport entity.Case) *cards.Card {
	_ = caseBuddyReport

	return cards.New(
		[]cards.Node{
			&cards.Container{
				Items: []cards.Node{
					&cards.TextBlock{
						Text: "Case Hygiene Alert",
					},
				},
			},
		}, []cards.Node{}).WithVersion(cards.Version12).WithSchema(cards.DefaultSchema)
}

func parseCaseBuddyAlert(alert string) ([]entity.Case, error) {
	doc, err := html.Parse(strings.NewReader(alert))
	if err != nil {
		return []entity.Case{}, err
	}

	var reports []entity.Case
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {
			var report entity.Case
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
						report.CaseNumber = data
					case 2:
						report.OwnerAlias = data
					case 3:
						report.CaseCreatedTime = data
					case 4:
						report.IncidentState = data
					case 5:
						report.CaseIdle = data
					case 6:
						report.CaseStatus = data
					}
					tdCount++

					slog.Info("Data: ", slog.String("data", data))
				}
			}
			if tdCount > 0 {
				reports = append(reports, report)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return reports, nil
}
