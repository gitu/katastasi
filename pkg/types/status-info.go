package types

type ServiceComponent struct {
	Status Status `json:"status"`
	Name   string `json:"name"`
	Info   string `json:"info"`
}

type Service struct {
	Id         string             `json:"id"`
	Name       string             `json:"name"`
	Status     Status             `json:"status"`
	LastUpdate int64              `json:"lastUpdate"`
	Components []ServiceComponent `json:"serviceComponents"`
}

type StatusInfo struct {
	LastUpdate    int64     `json:"lastUpdate"`
	OverallStatus Status    `json:"overallStatus"`
	Services      []Service `json:"services"`
	Name          string    `json:"name"`
}
