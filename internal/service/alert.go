package service

import (
	"bytes"
	"encoding/base64"
	"log/slog"
	"strings"

	"github.com/vermacodes/email-alert-processor/internal/entity"
	"golang.org/x/net/html"
)

type alertService struct{}

func NewAlertService() entity.AlertService {
	return &alertService{}
}

func (s *alertService) ProcessAlert(alert entity.Alert) (entity.Alert, error) {

	// Process alert here.
	// Alert message is base64 encoded HTML content. Decode it first.
	decodedAlertMessage, err := decodeBase64(alert.AlertMessage)
	if err != nil {
		slog.Error("Error decoding base64 string: ", err)
		return alert, err
	}

	// The alert is an HTML email, so we can use the HTML parser to extract the required information.
	texts, err := extractTextFromHTML(decodedAlertMessage)
	if err != nil {
		slog.Error("Error extracting text from HTML: ", err)
		return alert, err
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
	alert.AlertMessage = htmlDoc.String()

	return alert, nil
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

// func encodeBase64(s string) string {
// 	// Encode string to base64 here.
// 	return base64.StdEncoding.EncodeToString([]byte(s))
// }

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
