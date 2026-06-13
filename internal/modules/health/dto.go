package health

type StatusResponse struct {
	Status string            `json:"status"`
	Checks map[string]string `json:"checks"`
}
