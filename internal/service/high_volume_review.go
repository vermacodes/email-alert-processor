package service

import (
	"log/slog"
	"strings"

	cards "github.com/DanielTitkov/go-adaptive-cards"
	"github.com/vermacodes/email-alert-processor/internal/entity"
	"golang.org/x/net/html"
)

func highVolumeReview(alert entity.Alert) (entity.Alert, error) {
	slog.Info("Processing CaseHygiene alert")

	decodedAlertMessage, err := decodeBase64(alert.AlertMessage)
	if err != nil {
		slog.Error("Error decoding base64 string: ", err)
		return alert, err
	}

	caseBuddyReports, err := parseCaseBuddyAlertHighVolumeReview(decodedAlertMessage)
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

		adaptiveCard := buildAdaptiveCardHighVolumeReview(report)
		alert.TeamsMessages = append(alert.TeamsMessages, entity.TeamsMessage{
			AudType:      "user",
			AudId:        "ashisverma",
			AdaptiveCard: adaptiveCard,
		})
	}

	return alert, nil
}

func buildAdaptiveCardHighVolumeReview(caseBuddyReport entity.CaseBuddyReport) *cards.Card {
	return cards.New(
		[]cards.Node{
			&cards.TextBlock{
				Type:                cards.TextBlockType,
				Text:                "Weekly High Volume New Case (" + caseBuddyReport.CustomerName + ")",
				Size:                "large",
				Weight:              "lighter",
				HorizontalAlignment: "center",
			},
			&cards.Container{
				Type: "Container",
				Items: []cards.Node{
					&cards.ColumnSet{
						Type: cards.ColumnSetType,
						Columns: []*cards.Column{
							{
								Type: cards.ColumnType,
								Items: []cards.Node{
									&cards.ColumnSet{
										Type: cards.ColumnSetType,
										Columns: []*cards.Column{
											{
												Type: cards.ColumnType,
												Items: []cards.Node{
													&cards.TextBlock{
														Type:   cards.TextBlockType,
														Text:   "Case Number",
														Size:   "medium",
														Weight: "bolder",
													},
													&cards.TextBlock{
														Type:  cards.TextBlockType,
														Text:  caseBuddyReport.CaseNumber,
														Color: "accent",
													},
												},
											},
										},
										SelectAction: &cards.ActionOpenURL{
											Type: cards.ActionOpenURLType,
											URL:  caseBuddyReport.CaseLink,
										},
									},
								},
								Width: "auto",
							},
							{
								Type: cards.ColumnType,
								Items: []cards.Node{
									&cards.ColumnSet{
										Type: cards.ColumnSetType,
										Columns: []*cards.Column{
											{
												Type: cards.ColumnType,
												Items: []cards.Node{
													&cards.TextBlock{
														Type:   cards.TextBlockType,
														Text:   "Created",
														Size:   "medium",
														Weight: "bolder",
													},
													&cards.TextBlock{
														Type: cards.TextBlockType,
														Text: caseBuddyReport.CaseCreatedTime,
													},
												},
											},
										},
									},
								},
							},
							// {
							// 	Type: cards.ColumnType,
							// 	Items: []cards.Node{
							// 		&cards.ColumnSet{
							// 			Type: cards.ColumnSetType,
							// 			Columns: []*cards.Column{
							// 				{
							// 					Type: cards.ColumnType,
							// 					Items: []cards.Node{
							// 						&cards.TextBlock{
							// 							Type:   cards.TextBlockType,
							// 							Text:   "Customer",
							// 							Size:   "medium",
							// 							Weight: "bolder",
							// 						},
							// 						&cards.TextBlock{
							// 							Type: cards.TextBlockType,
							// 							Text: caseBuddyReport.CustomerName,
							// 						},
							// 					},
							// 				},
							// 			},
							// 		},
							// 	},
							// },
							{
								Type: cards.ColumnType,
								Items: []cards.Node{
									&cards.ColumnSet{
										Type: cards.ColumnSetType,
										Columns: []*cards.Column{
											{
												Type: cards.ColumnType,
												Items: []cards.Node{
													&cards.TextBlock{
														Type:   cards.TextBlockType,
														Text:   "Owner Alias",
														Size:   "medium",
														Weight: "bolder",
													},
													&cards.TextBlock{
														Type: cards.TextBlockType,
														Text: caseBuddyReport.OwnerAlias,
													},
												},
											},
										},
									},
								},
							},
							// {
							// 	Type: cards.ColumnType,
							// 	Items: []cards.Node{
							// 		&cards.ColumnSet{
							// 			Type: cards.ColumnSetType,
							// 			Columns: []*cards.Column{
							// 				{
							// 					Type: cards.ColumnType,
							// 					Items: []cards.Node{
							// 						&cards.TextBlock{
							// 							Type:   cards.TextBlockType,
							// 							Text:   "Owner Manager Alias",
							// 							Size:   "medium",
							// 							Weight: "bolder",
							// 						},
							// 						&cards.TextBlock{
							// 							Type: cards.TextBlockType,
							// 							Text: caseBuddyReport.OwnerManagerAlias,
							// 						},
							// 					},
							// 				},
							// 			},
							// 		},
							// 	},
							// },
						},
					},
				},
			},
		}, []cards.Node{}).WithSchema(cards.DefaultSchema).WithVersion(cards.Version12)
}

func parseCaseBuddyAlertHighVolumeReview(alert string) ([]entity.CaseBuddyReport, error) {
	doc, err := html.Parse(strings.NewReader(alert))
	if err != nil {
		return []entity.CaseBuddyReport{}, err
	}

	var reports []entity.CaseBuddyReport
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {
			var report entity.CaseBuddyReport
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
									report.CaseLink = attr.Val
									break
								}
							}
						} else {
							data = c.FirstChild.Data
						}
					}
					switch tdCount {
					case 1:
						report.CaseNumber = data
					case 2:
						report.CaseCreatedTime = data
					case 3:
						report.CustomerName = data
					case 4:
						report.OwnerAlias = data
					case 5:
						report.OwnerManagerAlias = data
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
