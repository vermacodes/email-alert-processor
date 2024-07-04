package entity

type Alert struct {
	ID           string `json:"id"`
	AlertType    string `json:"alertType"`
	AlertMessage string `json:"alertMessage"`
}

type AlertResponse struct {
	AlertID      string   `json:"alertId"`
	AlertType    string   `json:"alertType"`
	AlertMessage string   `json:"alertMessage"`
	Incident     Incident `json:"incident"`
}

type AlertService interface {
	ProcessAlert(alert Alert) ([]AlertResponse, error)
}
