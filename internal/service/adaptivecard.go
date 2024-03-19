package service

type Element struct {
	Type     string `json:"type"`
	Text     string `json:"text"`
	Weight   string `json:"weight"`
	Size     string `json:"size"`
	Wrap     bool   `json:"wrap"`
	IsSubtle bool   `json:"isSubtle"`
}

type AdaptiveCard struct {
	Type    string    `json:"type"`
	Version string    `json:"version"`
	Schema  string    `json:"$schema"`
	Body    []Element `json:"body"`
}

type AdaptiveCardService struct{}

func NewAdaptiveCardService() AdaptiveCardService {
	return AdaptiveCardService{}
}

func (s AdaptiveCardService) ProcessAdaptiveCard(card AdaptiveCard) (AdaptiveCard, error) {
	// Process adaptive card here.
	return card, nil
}
