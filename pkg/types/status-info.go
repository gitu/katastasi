package types

type Service struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Status     Status `json:"status"`
	LastUpdate int64  `json:"lastUpdate"`
}

type StatusInfo struct {
	LastUpdate    int64      `json:"lastUpdate"`
	StatusPage    StatusPage `json:"statusPage"`
	OverallStatus Status     `json:"overallStatus"`
	Services      []Service  `json:"services"`
}
