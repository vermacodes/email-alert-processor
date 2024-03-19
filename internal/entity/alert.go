package entity

type Alert struct {
	ID           string `json:"id"`
	AlertType    string `json:"alertType"`
	AlertMessage string `json:"alertMessage"`
}

type ALertResponse struct {
	AlertID      string   `json:"alertId"`
	AlertType    string   `json:"alertType"`
	AlertMessage string   `json:"alertMessage"`
	Incident     Incident `json:"incident"`
}

type AlertService interface {
	ProcessAlert(alert Alert) (Alert, error)
}
