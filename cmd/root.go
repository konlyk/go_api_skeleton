package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/konlyk/go_api_skeleton/bootstrap"
)

func execute() error {
	return newRootCmd().ExecuteContext(context.Background())
}

func newRootCmd() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:          "go-api-skeleton",
		Short:        "Run the Go API skeleton service",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, stop := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()
			return run(ctx, configPath)
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "", "Path to config YAML file (defaults to ./config.yaml when present)")
	return cmd
}

func run(ctx context.Context, configPath string) error {
	app, err := bootstrap.NewApplicationWithConfig(ctx, configPath)
	if err != nil {
		return fmt.Errorf("initialize application: %w", err)
	}

	engine, err := bootstrap.SetupRouter(app)
	if err != nil {
		return fmt.Errorf("setup routes: %w", err)
	}

	server := &http.Server{
		Addr:              app.Config.HTTPAddress(),
		Handler:           engine,
		ReadTimeout:       app.Config.HTTPReadTimeout,
		ReadHeaderTimeout: app.Config.HTTPHeaderTimeout,
		WriteTimeout:      app.Config.HTTPWriteTimeout,
		IdleTimeout:       app.Config.HTTPIdleTimeout,
	}

	errCh := make(chan error, 1)
	go func() {
		app.Observability.Logger.Info().Str("addr", app.Config.HTTPAddress()).Msg("http server starting")
		if serveErr := server.ListenAndServe(); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			errCh <- fmt.Errorf("listen and serve: %w", serveErr)
		}
		close(errCh)
	}()

	select {
	case serveErr := <-errCh:
		if serveErr != nil {
			shutdownCtx, cancel := context.WithTimeout(context.Background(), app.Config.ShutdownTimeout)
			defer cancel()
			if shutdownErr := app.Shutdown(shutdownCtx); shutdownErr != nil {
				app.Observability.Logger.Error().Err(shutdownErr).Msg("shutdown observability")
			}
			return serveErr
		}
		return nil
	case <-ctx.Done():
		app.Observability.Logger.Info().Msg("shutdown signal received")
	}

	app.Health.SetReady(false)
	if app.Config.ShutdownDrainDelay > 0 {
		timer := time.NewTimer(app.Config.ShutdownDrainDelay)
		defer timer.Stop()
		select {
		case <-timer.C:
		case serveErr := <-errCh:
			if serveErr != nil {
				return serveErr
			}
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), app.Config.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		if closeErr := server.Close(); closeErr != nil {
			app.Observability.Logger.Error().Err(closeErr).Msg("force close server")
		}
		return fmt.Errorf("shutdown server: %w", err)
	}

	if err := app.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown observability: %w", err)
	}

	app.Observability.Logger.Info().Msg("shutdown complete")
	return nil
}
