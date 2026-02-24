package bootstrap

import (
	"context"

	"github.com/konlyk/go_api_skeleton/controller"
	"github.com/konlyk/go_api_skeleton/domain"
)

type Application struct {
	Config        *domain.Config
	Health        *HealthManager
	Observability *Observability
	Controllers   *controller.Services
}

func NewApplication(ctx context.Context) (*Application, error) {
	return NewApplicationWithConfig(ctx, "")
}

func NewApplicationWithConfig(ctx context.Context, configPath string) (*Application, error) {
	config, err := domain.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	observability, err := NewObservability(ctx, config)
	if err != nil {
		return nil, err
	}

	health := NewHealthManager(config.ServiceName)
	controllers := controller.Setup(health)

	return &Application{
		Config:        config,
		Health:        health,
		Observability: observability,
		Controllers:   controllers,
	}, nil
}

func (a *Application) Shutdown(ctx context.Context) error {
	return a.Observability.Shutdown(ctx)
}
