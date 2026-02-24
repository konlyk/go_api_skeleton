package bootstrap

import (
	"sync/atomic"

	"github.com/konlyk/go_api_skeleton/domain"
)

type HealthManager struct {
	service string
	ready   atomic.Bool
}

func NewHealthManager(service string) *HealthManager {
	h := &HealthManager{service: service}
	h.ready.Store(true)
	return h
}

func (h *HealthManager) SetReady(ready bool) {
	h.ready.Store(ready)
}

func (h *HealthManager) Liveness() domain.HealthStatus {
	return domain.HealthStatus{Status: "ok", Service: h.service}
}

func (h *HealthManager) Readiness() domain.HealthStatus {
	if h.ready.Load() {
		return domain.HealthStatus{Status: "ok", Service: h.service}
	}
	return domain.HealthStatus{Status: "not_ready", Service: h.service}
}

func (h *HealthManager) IsReady() bool {
	return h.ready.Load()
}
