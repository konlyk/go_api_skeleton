package domain

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceName        string        `mapstructure:"service_name"`
	Port               string        `mapstructure:"port"`
	PrivateAPIToken    string        `mapstructure:"private_api_token"`
	HTTPReadTimeout    time.Duration `mapstructure:"http_read_timeout"`
	HTTPHeaderTimeout  time.Duration `mapstructure:"http_read_header_timeout"`
	HTTPWriteTimeout   time.Duration `mapstructure:"http_write_timeout"`
	HTTPIdleTimeout    time.Duration `mapstructure:"http_idle_timeout"`
	ShutdownDrainDelay time.Duration `mapstructure:"shutdown_drain_delay"`
	ShutdownTimeout    time.Duration `mapstructure:"shutdown_timeout"`
	LogLevel           zerolog.Level `mapstructure:"log_level"`
	EnableTracing      bool          `mapstructure:"enable_otel_tracing"`
	TraceSampleRatio   float64       `mapstructure:"trace_sample_ratio"`
}

func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()
	setConfigDefaults(v)

	if err := bindConfigEnv(v); err != nil {
		return nil, err
	}

	if err := readConfigFile(v, configPath); err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg, viper.DecodeHook(
		mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			stringToZerologLevelHookFunc(),
		),
	)); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}

	if cfg.TraceSampleRatio < 0 || cfg.TraceSampleRatio > 1 {
		return nil, fmt.Errorf("TRACE_SAMPLE_RATIO must be between 0 and 1")
	}

	return cfg, nil
}

func (c *Config) HTTPAddress() string {
	return ":" + c.Port
}

func setConfigDefaults(v *viper.Viper) {
	v.SetDefault("service_name", "go-api-skeleton")
	v.SetDefault("port", "8080")
	v.SetDefault("private_api_token", "dev-private-token")
	v.SetDefault("http_read_timeout", "15s")
	v.SetDefault("http_read_header_timeout", "5s")
	v.SetDefault("http_write_timeout", "15s")
	v.SetDefault("http_idle_timeout", "60s")
	v.SetDefault("shutdown_drain_delay", "5s")
	v.SetDefault("shutdown_timeout", "20s")
	v.SetDefault("log_level", "info")
	v.SetDefault("enable_otel_tracing", false)
	v.SetDefault("trace_sample_ratio", 0.1)
}

func bindConfigEnv(v *viper.Viper) error {
	bindings := map[string]string{
		"service_name":             "SERVICE_NAME",
		"port":                     "PORT",
		"private_api_token":        "PRIVATE_API_TOKEN",
		"http_read_timeout":        "HTTP_READ_TIMEOUT",
		"http_read_header_timeout": "HTTP_READ_HEADER_TIMEOUT",
		"http_write_timeout":       "HTTP_WRITE_TIMEOUT",
		"http_idle_timeout":        "HTTP_IDLE_TIMEOUT",
		"shutdown_drain_delay":     "SHUTDOWN_DRAIN_DELAY",
		"shutdown_timeout":         "SHUTDOWN_TIMEOUT",
		"log_level":                "LOG_LEVEL",
		"enable_otel_tracing":      "ENABLE_OTEL_TRACING",
		"trace_sample_ratio":       "TRACE_SAMPLE_RATIO",
	}

	for key, envName := range bindings {
		if err := v.BindEnv(key, envName); err != nil {
			return fmt.Errorf("bind env %s: %w", envName, err)
		}
	}

	v.AutomaticEnv()
	return nil
}

func readConfigFile(v *viper.Viper, configPath string) error {
	if strings.TrimSpace(configPath) != "" {
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			return fmt.Errorf("read config file %s: %w", configPath, err)
		}
		return nil
	}

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		var configNotFound viper.ConfigFileNotFoundError
		if errors.As(err, &configNotFound) {
			return nil
		}
		return fmt.Errorf("read config.yaml: %w", err)
	}

	return nil
}

func stringToZerologLevelHookFunc() mapstructure.DecodeHookFuncType {
	logLevelType := reflect.TypeOf(zerolog.InfoLevel)

	return func(from reflect.Type, to reflect.Type, data any) (any, error) {
		if from.Kind() != reflect.String || to != logLevelType {
			return data, nil
		}

		level, err := zerolog.ParseLevel(strings.ToLower(strings.TrimSpace(data.(string))))
		if err != nil {
			return nil, fmt.Errorf("unsupported LOG_LEVEL %q (expected debug|info|warn|error)", data)
		}

		return level, nil
	}
}
