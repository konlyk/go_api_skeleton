package bootstrap

import (
	"context"
	"fmt"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/konlyk/go_api_skeleton/domain"
)

type Observability struct {
	Logger          zerolog.Logger
	MetricsRegistry *prometheus.Registry
	Tracer          trace.Tracer
	shutdown        func(context.Context) error
}

func NewObservability(ctx context.Context, config *domain.Config) (*Observability, error) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("service", config.ServiceName).Logger().Level(config.LogLevel)

	registry := prometheus.NewRegistry()
	if err := registry.Register(collectors.NewGoCollector()); err != nil {
		return nil, fmt.Errorf("register go collector: %w", err)
	}
	if err := registry.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{})); err != nil {
		return nil, fmt.Errorf("register process collector: %w", err)
	}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	obs := &Observability{
		Logger:          logger,
		MetricsRegistry: registry,
		Tracer:          otel.Tracer(config.ServiceName),
		shutdown:        func(context.Context) error { return nil },
	}

	if !config.EnableTracing {
		return obs, nil
	}

	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("create otlp trace exporter: %w", err)
	}

	res, err := resource.New(ctx, resource.WithAttributes(semconv.ServiceName(config.ServiceName)))
	if err != nil {
		return nil, fmt.Errorf("create otel resource: %w", err)
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(config.TraceSampleRatio)),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(provider)

	obs.Tracer = provider.Tracer(config.ServiceName)
	obs.shutdown = provider.Shutdown

	return obs, nil
}

func (o *Observability) Shutdown(ctx context.Context) error {
	return o.shutdown(ctx)
}
