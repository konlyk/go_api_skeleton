package controller

import (
	"context"

	"github.com/konlyk/go_api_skeleton/api/openapi"
	"github.com/konlyk/go_api_skeleton/domain"
)

type HelloController struct {
	helloUsecase domain.HelloUsecase
}

func NewHelloController(helloUsecase domain.HelloUsecase) *HelloController {
	return &HelloController{helloUsecase: helloUsecase}
}

func (c *HelloController) GetHello(ctx context.Context, _ openapi.GetHelloRequestObject) (openapi.GetHelloResponseObject, error) {
	hello, err := c.helloUsecase.Execute(ctx)
	if err != nil {
		return openapi.GetHello500JSONResponse{Error: "internal server error"}, nil
	}

	return openapi.GetHello200JSONResponse{
		Message:   hello.Message,
		Timestamp: hello.Timestamp,
	}, nil
}
