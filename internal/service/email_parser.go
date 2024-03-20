package service

import (
	"log/slog"
	"strings"

	"golang.org/x/net/html"

	"github.com/vermacodes/email-alert-processor/internal/entity"
)

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
