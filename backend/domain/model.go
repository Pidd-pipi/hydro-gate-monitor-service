package domain

type Gate struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Basin           string `json:"basin"`
	State           string `json:"state"`
	LastObservation string `json:"lastObservation"`
	AlertLevel      string `json:"alertLevel"`
	Acknowledged    bool   `json:"acknowledged"`
}

type AcknowledgeRequest struct {
	Note string `json:"note"`
}
