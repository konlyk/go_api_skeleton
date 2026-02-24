package controller

import (
	"github.com/konlyk/go_api_skeleton/api/openapi"
	"github.com/konlyk/go_api_skeleton/domain"
	"github.com/konlyk/go_api_skeleton/repository"
	"github.com/konlyk/go_api_skeleton/usecase"
)

var _ openapi.StrictServerInterface = (*Services)(nil)

type Services struct {
	*PingController
	*HelloController
	*HealthController
}

func Setup(health domain.HealthState) *Services {
	clockRepository := repository.NewSystemClockRepository()

	pingUsecase := usecase.NewPingUsecase(clockRepository)
	helloUsecase := usecase.NewHelloUsecase(clockRepository)

	return &Services{
		PingController:   NewPingController(pingUsecase),
		HelloController:  NewHelloController(helloUsecase),
		HealthController: NewHealthController(health),
	}
}
