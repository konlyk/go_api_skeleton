package domain

import (
	"context"
	"time"
)

type Hello struct {
	Message   string
	Timestamp time.Time
}

type HelloUsecase interface {
	Execute(ctx context.Context) (Hello, error)
}
