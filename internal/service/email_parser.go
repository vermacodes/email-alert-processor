package service

import (
	"bytes"
	"log/slog"
	"strings"

	"golang.org/x/net/html"

	"github.com/vermacodes/email-alert-processor/internal/entity"
)

func parseCaseBuddyEmail(alert string) ([]entity.Incident, error) {
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

					if c.FirstChild == nil {
						tdCount++
						continue
					}

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
						incident.IncidentCreatedTime = data
					case 3:
						incident.IncidentCustomerName = data
					case 4:
						incident.IncidentOwnerAlias = data
					case 5:
						incident.IncidentOwnerManagerAlias = data
					case 6:
						incident.IncidentState = data
					case 7:
						incident.IncidentIdleDuration = data
					case 8:
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

	// Get all the nodes of type p.
	nodes, err := getNodesOfGivenDataType(doc, "p")
	if err != nil {
		return nil, err
	}

	var texts []string

	for _, node := range nodes {
		var buf bytes.Buffer
		err := html.Render(&buf, node.FirstChild)
		if err != nil {
			return nil, err
		}

		text, err := getInnerMostTextFromFirstChild(node)
		if err != nil {
			return nil, err
		}

		texts = append(texts, text)

		slog.Info("Extracted text", slog.String("text", text))
	}

	return texts, nil
}

func parseTextFromIridiasEmail(extractedEmailText []string) (entity.Incident, error) {
	var incident entity.Incident
	filters := []string{"Ticket", "Severity", "Status", "Customer", "Product", "Created On"}
	for i := 0; i < len(extractedEmailText); i++ {
		// Check if the current element is in the filters slice.
		isFilter := false
		filterValue := ""
		for _, filter := range filters {
			if strings.Contains(extractedEmailText[i], filter) {
				slog.Info("Filter: ", slog.String("filter", filter))
				isFilter = true
				filterValue = filter
				break
			}
		}

		// If the current element is in the filters slice, add it and the next one to the HTML document.
		if isFilter && i+1 < len(extractedEmailText) {
			switch filterValue {
			case "Ticket":
				slog.Info("Ticket: ", slog.String("ticket", extractedEmailText[i+1]))
				incident.IncidentNumber = extractedEmailText[i+1]
			case "Severity":
				incident.IncidentSeverity = extractedEmailText[i+1]
			case "Status":
				incident.IncidentStatus = extractedEmailText[i+1]
			case "Customer":
				incident.IncidentCustomerName = extractedEmailText[i+1]
			case "Product":
				slog.Info("Product: ", slog.String("product", extractedEmailText[i+1]))
				incident.IncidentProduct = extractedEmailText[i+1]
			case "Created On":
				incident.IncidentCreatedTime = extractedEmailText[i+1]
			}
			i++ // Skip the next element because it has already been added.
		}
	}
	return incident, nil
}

func getNodesOfGivenDataType(doc *html.Node, dataType string) ([]*html.Node, error) {
	// Loop through all nodes and find the nodes of given data type.
	var nodes []*html.Node

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == dataType {
			nodes = append(nodes, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	return nodes, nil
}

func getInnerMostTextFromFirstChild(doc *html.Node) (string, error) {
	// Loop throught nodes and find the inner most text from the first child.
	// If the first child is a text node, and got no children. return the data.

	var f func(*html.Node) string
	f = func(n *html.Node) string {
		if n == nil {
			return ""
		}

		if n.Type == html.TextNode && n.FirstChild == nil {
			return n.Data
		}
		return f(n.FirstChild)
	}

	return f(doc), nil
}
