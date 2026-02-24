package domain

import (
	"context"
	"time"
)

type Ping struct {
	Message   string
	Timestamp time.Time
}

type PingUsecase interface {
	Execute(ctx context.Context) (Ping, error)
}

type ClockRepository interface {
	Now() time.Time
}
