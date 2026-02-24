package controller

import (
	"context"

	"github.com/konlyk/go_api_skeleton/api/openapi"
	"github.com/konlyk/go_api_skeleton/domain"
)

type PingController struct {
	pingUsecase domain.PingUsecase
}

func NewPingController(pingUsecase domain.PingUsecase) *PingController {
	return &PingController{pingUsecase: pingUsecase}
}

func (c *PingController) GetPing(ctx context.Context, _ openapi.GetPingRequestObject) (openapi.GetPingResponseObject, error) {
	ping, err := c.pingUsecase.Execute(ctx)
	if err != nil {
		return openapi.GetPing500JSONResponse{Error: "internal server error"}, nil
	}

	return openapi.GetPing200JSONResponse{
		Message:   ping.Message,
		Timestamp: ping.Timestamp,
	}, nil
}
