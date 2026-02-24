package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/konlyk/go_api_skeleton/domain"
)

type fakeClock struct {
	now time.Time
}

func (f fakeClock) Now() time.Time { return f.now }

var _ domain.ClockRepository = (*fakeClock)(nil)

func TestHelloUsecaseExecute(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.January, 2, 10, 0, 0, 0, time.UTC)
	uc := NewHelloUsecase(fakeClock{now: now})

	hello, err := uc.Execute(context.Background())
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if hello.Message != "hello world" {
		t.Fatalf("expected message hello world, got %q", hello.Message)
	}
	if !hello.Timestamp.Equal(now) {
		t.Fatalf("expected timestamp %s, got %s", now, hello.Timestamp)
	}
}
