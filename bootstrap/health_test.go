package bootstrap

import "testing"

func TestHealthManagerTransitions(t *testing.T) {
	t.Parallel()

	health := NewHealthManager("svc")
	if !health.IsReady() {
		t.Fatal("expected service to be ready by default")
	}
	if health.Readiness().Status != "ok" {
		t.Fatalf("expected readiness status ok, got %q", health.Readiness().Status)
	}

	health.SetReady(false)
	if health.IsReady() {
		t.Fatal("expected service to be not ready")
	}
	if health.Readiness().Status != "not_ready" {
		t.Fatalf("expected readiness status not_ready, got %q", health.Readiness().Status)
	}
}
