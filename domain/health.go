package domain

type HealthStatus struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

type HealthState interface {
	Liveness() HealthStatus
	Readiness() HealthStatus
	IsReady() bool
}
