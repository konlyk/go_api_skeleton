package usecase

import (
	"context"

	"github.com/konlyk/go_api_skeleton/domain"
)

type helloUsecase struct {
	clock domain.ClockRepository
}

func NewHelloUsecase(clock domain.ClockRepository) domain.HelloUsecase {
	return &helloUsecase{clock: clock}
}

func (u *helloUsecase) Execute(context.Context) (domain.Hello, error) {
	return domain.Hello{Message: "hello world", Timestamp: u.clock.Now()}, nil
}
