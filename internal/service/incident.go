package service

import (
	"log/slog"

	"github.com/vermacodes/email-alert-processor/internal/entity"
)

func parseCaseBuddyEmail(encodedEmailBody string) ([]entity.Case, error) {
  // Decode the base64 encoded email encodedEmailBody
  decodedEmailBody, err := decodeBase64(encodedEmailBody)
  if err != nil {
    slog.Error("Error decoding base64 string: ", err)
    return []entity.Case{}, err
  }
  // The email is an HTML email, so we can use the HTML parser to extract the required information.
	doc, err := html.parse(strings.newreader(decodedEmailBody))
	if err != nil {
		return []entity.Case{}, err
	}

	var cases []entity.Case
	var f func(*html.node)
	f = func(n *html.node) {
		if n.type == html.elementnode && n.data == "tr" {
			var case entity.Case
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
									case.CaseLink = attr.Val
									break
								}
							}
						} else {
							data = c.FirstChild.Data
						}
					}
					switch tdCount {
					case 1:
						case.CaseNumber = data
					case 2:
						case.CaseCreatedTime = data
					case 3:
						case.CustomerName = data
					case 4:
						case.OwnerAlias = data
					case 5:
						case.OwnerManagerAlias = data
					}
					tdCount++

					slog.Info("Data: ", slog.String("data", data))
				}
			}
			if tdCount > 0 {
				cases = append(cases, case)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return cases, nil
}
