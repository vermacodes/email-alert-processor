package service

import (
	"log/slog"
	"strings"

	"golang.org/x/net/html"

	"github.com/vermacodes/email-alert-processor/internal/entity"
)

func parseCaseBuddyEmail(encodedEmailBody string) ([]entity.Incident, error) {
	// Decode the base64 encoded email encodedEmailBody
	decodedEmailBody, err := decodeBase64(encodedEmailBody)
	if err != nil {
		slog.Error("Error decoding base64 string: ", err)
		return []entity.Incident{}, err
	}
	// The email is an HTML email, so we can use the HTML parser to extract the required information.
	doc, err := html.Parse(strings.NewReader(decodedEmailBody))
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
					if c.FirstChild != nil {
						if c.FirstChild.Data == "a" {
							if c.FirstChild.FirstChild != nil {
								data = c.FirstChild.FirstChild.Data
							}

							// Get the case link
							for _, attr := range c.FirstChild.Attr {
								if attr.Key == "href" {
									incident.IncidentLink = attr.Val
									break
								}
							}
						} else {
							data = c.FirstChild.Data
						}
					}
					switch tdCount {
					case 1:
						incident.IncidentNumber = data
					case 2:
						incident.IncidentCreatedTime = data
					case 3:
						incident.IncidentCustomerName = data
					case 4:
						incident.IncidentOwnerAlias = data
					case 5:
						incident.IncidentOwnerManagerAlias = data
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
