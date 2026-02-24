package usecase

import (
	"context"

	"github.com/konlyk/go_api_skeleton/domain"
)

type pingUsecase struct {
	clock domain.ClockRepository
}

func NewPingUsecase(clock domain.ClockRepository) domain.PingUsecase {
	return &pingUsecase{clock: clock}
}

func (u *pingUsecase) Execute(context.Context) (domain.Ping, error) {
	return domain.Ping{Message: "pong", Timestamp: u.clock.Now()}, nil
}
