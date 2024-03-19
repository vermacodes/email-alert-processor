package entity

type Incident struct {
	IncidentNumber            string `json:"incidentNumber"`
	IncidentLink              string `json:"incidentLink"`
	IncidentOwnerAlias        string `json:"incidentOwnerAlias"`
	IncidentOwnerManagerAlias string `json:"incidentOwnerManagerAlias"`
	IncidentCreatedTime       string `json:"incidentCreatedTime"`
	IncidentState             string `json:"incidentState"`
	IncidentIdleDuration      string `json:"incidentIdleDuration"`
	IncidentStatus            string `json:"incidentStatus"`
	IncidentCustomerName      string `json:"incidentCustomerName"`
}
