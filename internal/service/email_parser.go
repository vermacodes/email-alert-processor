package service

import (
	"bytes"
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

func extractTextFromIridiasEmail(htmlContent string) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}
	var texts []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "p" {
			var buf bytes.Buffer
			err := html.Render(&buf, n)
			if err != nil {
				return
			}
			texts = append(texts, buf.String())
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return texts, nil
}

func parseTextFromIridiasEmail(extractedEmailText []string) (entity.Incident, error) {
	var incident entity.Incident
	filters := []string{"Ticket", "Severity", "Status", "Customer", "Product", "Created On"}
	for i := 0; i < len(extractedEmailText); i++ {
		// Check if the current element is in the filters slice.
		isFilter := false
		for _, filter := range filters {
			if strings.Contains(extractedEmailText[i], filter) {
				isFilter = true
				break
			}
		}

		// If the current element is in the filters slice, add it and the next one to the HTML document.
		if isFilter && i+1 < len(extractedEmailText) {
			switch extractedEmailText[i] {
			case "Ticket":
				incident.IncidentNumber = extractedEmailText[i+1]
			case "Severity":
				incident.IncidentSeverity = extractedEmailText[i+1]
			case "Status":
				incident.IncidentStatus = extractedEmailText[i+1]
			case "Customer":
				incident.IncidentCustomerName = extractedEmailText[i+1]
			case "Product":
				incident.IncidentProduct = extractedEmailText[i+1]
			case "Created On":
				incident.IncidentCreatedTime = extractedEmailText[i+1]
			}
			i++ // Skip the next element because it has already been added.
		}
	}
	return incident, nil
}
