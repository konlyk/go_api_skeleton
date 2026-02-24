package controller

import (
	"context"

	"github.com/konlyk/go_api_skeleton/api/openapi"
	"github.com/konlyk/go_api_skeleton/domain"
)

type HealthController struct {
	health domain.HealthState
}

func NewHealthController(health domain.HealthState) *HealthController {
	return &HealthController{health: health}
}

func (c *HealthController) GetLivez(_ context.Context, _ openapi.GetLivezRequestObject) (openapi.GetLivezResponseObject, error) {
	status := c.health.Liveness()
	return openapi.GetLivez200JSONResponse(toHealthResponse(status)), nil
}

func (c *HealthController) GetReadyz(_ context.Context, _ openapi.GetReadyzRequestObject) (openapi.GetReadyzResponseObject, error) {
	status := c.health.Readiness()
	if c.health.IsReady() {
		return openapi.GetReadyz200JSONResponse(toHealthResponse(status)), nil
	}
	return openapi.GetReadyz503JSONResponse(toHealthResponse(status)), nil
}

func toHealthResponse(status domain.HealthStatus) openapi.HealthResponse {
	value := openapi.NotReady
	if status.Status == "ok" {
		value = openapi.Ok
	}

	return openapi.HealthResponse{Status: value, Service: status.Service}
}
