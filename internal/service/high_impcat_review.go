package service

import (
	"log/slog"
	"strings"

	cards "github.com/DanielTitkov/go-adaptive-cards"
	"github.com/vermacodes/email-alert-processor/internal/entity"
	"golang.org/x/net/html"
)

func highImpactReview(alert entity.Alert) (entity.Alert, error) {
	slog.Info("Processing CaseHygiene alert")

	decodedAlertMessage, err := decodeBase64(alert.AlertMessage)
	if err != nil {
		slog.Error("Error decoding base64 string: ", err)
		return alert, err
	}

	caseBuddyReports, err := parseCaseBuddyAlertHighImpactReview(decodedAlertMessage)
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

		adaptiveCard := buildAdaptiveCardHighImpactReview(report)
		alert.TeamsMessages = append(alert.TeamsMessages, entity.TeamsMessage{
			AudType:      "user",
			AudId:        "ashisverma",
			AdaptiveCard: adaptiveCard,
		})
	}

	return alert, nil
}

func buildadaptivecardhighimpactreview(casebuddyreport entity.casebuddyreport) *cards.card {
	return cards.new(
		[]cards.node{
			&cards.textblock{
				type:                cards.textblocktype,
				text:                "high impact new case (" + casebuddyreport.customername + ")",
				size:                "large",
				weight:              "lighter",
				horizontalalignment: "center",
			},
			&cards.container{
				type: "container",
				items: []cards.node{
					&cards.columnset{
						type: cards.columnsettype,
						columns: []*cards.column{
							{
								type: cards.columntype,
								items: []cards.node{
									&cards.columnset{
										type: cards.columnsettype,
										columns: []*cards.column{
											{
												type: cards.columntype,
												items: []cards.node{
													&cards.textblock{
														type:   cards.textblocktype,
														text:   "case number",
														size:   "medium",
														weight: "bolder",
													},
													&cards.textblock{
														type:  cards.textblocktype,
														text:  casebuddyreport.casenumber,
														color: "accent",
													},
												},
											},
										},
										selectaction: &cards.actionopenurl{
											type: cards.actionopenurltype,
											url:  casebuddyreport.caselink,
										},
									},
								},
								width: "auto",
							},
							{
								type: cards.columntype,
								items: []cards.node{
									&cards.columnset{
										type: cards.columnsettype,
										columns: []*cards.column{
											{
												type: cards.columntype,
												items: []cards.node{
													&cards.textblock{
														type:   cards.textblocktype,
														text:   "created",
														size:   "medium",
														weight: "bolder",
													},
													&cards.textblock{
														type: cards.textblocktype,
														text: casebuddyreport.casecreatedtime,
													},
												},
											},
										},
									},
								},
							},
							// {
							// 	type: cards.columntype,
							// 	items: []cards.node{
							// 		&cards.columnset{
							// 			type: cards.columnsettype,
							// 			columns: []*cards.column{
							// 				{
							// 					type: cards.columntype,
							// 					items: []cards.node{
							// 						&cards.textblock{
							// 							type:   cards.textblocktype,
							// 							text:   "customer",
							// 							size:   "medium",
							// 							weight: "bolder",
							// 						},
							// 						&cards.textblock{
							// 							type: cards.textblocktype,
							// 							text: casebuddyreport.customername,
							// 						},
							// 					},
							// 				},
							// 			},
							// 		},
							// 	},
							// },
							{
								type: cards.columntype,
								items: []cards.node{
									&cards.columnset{
										type: cards.columnsettype,
										columns: []*cards.column{
											{
												type: cards.columntype,
												items: []cards.node{
													&cards.textblock{
														type:   cards.textblocktype,
														text:   "owner alias",
														size:   "medium",
														weight: "bolder",
													},
													&cards.textblock{
														type: cards.textblocktype,
														text: casebuddyreport.owneralias,
													},
												},
											},
										},
									},
								},
							},
							// {
							// 	type: cards.columntype,
							// 	items: []cards.node{
							// 		&cards.columnset{
							// 			type: cards.columnsettype,
							// 			columns: []*cards.column{
							// 				{
							// 					type: cards.columntype,
							// 					items: []cards.node{
							// 						&cards.textblock{
							// 							type:   cards.textblocktype,
							// 							text:   "owner manager alias",
							// 							size:   "medium",
							// 							weight: "bolder",
							// 						},
							// 						&cards.textblock{
							// 							type: cards.textblocktype,
							// 							text: casebuddyreport.ownermanageralias,
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
		}, []cards.node{}).withschema(cards.defaultschema).withversion(cards.version12)
}

func parsecasebuddyalerthighimpactreview(alert string) ([]entity.casebuddyreport, error) {
	doc, err := html.parse(strings.newreader(alert))
	if err != nil {
		return []entity.casebuddyreport{}, err
	}

	var reports []entity.casebuddyreport
	var f func(*html.node)
	f = func(n *html.node) {
		if n.type == html.elementnode && n.data == "tr" {
			var report entity.casebuddyreport
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
