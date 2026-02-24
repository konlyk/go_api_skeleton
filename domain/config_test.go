package domain

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigDefaults(t *testing.T) {
	t.Setenv("SERVICE_NAME", "svc")
	t.Setenv("PORT", "9090")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("HTTP_READ_TIMEOUT", "11s")
	t.Setenv("TRACE_SAMPLE_RATIO", "0.25")

	cfg, err := LoadConfig("")
	if err != nil {
		t.Fatalf("expected config to load, got error: %v", err)
	}

	if cfg.ServiceName != "svc" {
		t.Fatalf("expected service name svc, got %q", cfg.ServiceName)
	}
	if cfg.Port != "9090" {
		t.Fatalf("expected port 9090, got %q", cfg.Port)
	}
	if cfg.LogLevel.String() != "debug" {
		t.Fatalf("expected log level debug, got %q", cfg.LogLevel.String())
	}
	if cfg.HTTPReadTimeout.String() != "11s" {
		t.Fatalf("expected read timeout 11s, got %s", cfg.HTTPReadTimeout)
	}
	if cfg.TraceSampleRatio != 0.25 {
		t.Fatalf("expected trace ratio 0.25, got %v", cfg.TraceSampleRatio)
	}
}

func TestLoadConfigInvalidLogLevel(t *testing.T) {
	t.Setenv("LOG_LEVEL", "invalid-level")

	if _, err := LoadConfig(""); err == nil {
		t.Fatal("expected invalid LOG_LEVEL to fail")
	}
}

func TestLoadConfigFromYAML(t *testing.T) {
	configFilePath := filepath.Join(t.TempDir(), "config.yaml")
	configYAML := []byte("service_name: from-file\nport: \"7777\"\nlog_level: warn\n")
	if err := os.WriteFile(configFilePath, configYAML, 0o644); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	cfg, err := LoadConfig(configFilePath)
	if err != nil {
		t.Fatalf("expected yaml config to load, got error: %v", err)
	}

	if cfg.ServiceName != "from-file" {
		t.Fatalf("expected service name from-file, got %q", cfg.ServiceName)
	}
	if cfg.Port != "7777" {
		t.Fatalf("expected port 7777, got %q", cfg.Port)
	}
	if cfg.LogLevel.String() != "warn" {
		t.Fatalf("expected log level warn, got %q", cfg.LogLevel.String())
	}
}

func TestLoadConfigEnvOverridesYAML(t *testing.T) {
	configFilePath := filepath.Join(t.TempDir(), "config.yaml")
	configYAML := []byte("service_name: from-file\nport: \"7777\"\n")
	if err := os.WriteFile(configFilePath, configYAML, 0o644); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	t.Setenv("PORT", "8888")

	cfg, err := LoadConfig(configFilePath)
	if err != nil {
		t.Fatalf("expected config to load, got error: %v", err)
	}

	if cfg.Port != "8888" {
		t.Fatalf("expected env to override file port to 8888, got %q", cfg.Port)
	}
}
