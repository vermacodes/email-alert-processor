package service

import (
	"bytes"
	"encoding/base64"
	"log/slog"
	"strings"

	"golang.org/x/net/html"

	"github.com/vermacodes/email-alert-processor/internal/entity"
)

func cnflntNewCaseAlert(alert entity.Alert) ([]entity.AlertResponse, error) {
	// Process alert here.
	// Alert message is base64 encoded HTML content. Decode it first.
	decodedAlertMessage, err := decodeBase64(alert.AlertMessage)
	if err != nil {
		slog.Error("Error decoding base64 string: ", err)
		return []entity.AlertResponse{}, err
	}

	// The alert is an HTML email, so we can use the HTML parser to extract the required information.
	texts, err := extractTextFromHTML(decodedAlertMessage)
	if err != nil {
		slog.Error("Error extracting text from HTML: ", err)
		return []entity.AlertResponse{}, err
	}

	// Now we can use the extracted texts to process the alert.
	var htmlDoc strings.Builder

	htmlDoc.WriteString("<html><body>")

	filters := []string{"Ticket", "Severity", "Status", "Customer", "Product", "Created On"}

	for i := 0; i < len(texts); i++ {
		// Check if the current element is in the filters slice.
		isFilter := false
		for _, filter := range filters {
			if strings.Contains(texts[i], filter) {
				isFilter = true
				break
			}
		}

		// If the current element is in the filters slice, add it and the next one to the HTML document.
		if isFilter && i+1 < len(texts) {
			slog.Info("Processing text: ", slog.String("text", texts[i]))
			htmlDoc.WriteString(texts[i])
			htmlDoc.WriteString(texts[i+1])
			i++ // Skip the next element because it has already been added.
		}
	}

	htmlDoc.WriteString("</body></html>")

	// Now htmlDoc.String() contains the HTML document.

	alertResponse := entity.AlertResponse{
		AlertID:      alert.ID,
		AlertType:    alert.AlertType,
		AlertMessage: htmlDoc.String(),
	}

	alertResponses := []entity.AlertResponse{alertResponse}

	return alertResponses, nil
}

func decodeBase64(s string) (string, error) {
	// Decode base64 encoded string here.
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	// Convert byte array to string.
	return string(data), nil
}

func extractTextFromHTML(htmlContent string) ([]string, error) {
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
