package entity

type Alert struct {
	ID           string `json:"id"`
	AlertType    string `json:"alertType"`
	AlertMessage string `json:"alertMessage"`
}

type AlertService interface {
	ProcessAlert(alert Alert) (Alert, error)
}
