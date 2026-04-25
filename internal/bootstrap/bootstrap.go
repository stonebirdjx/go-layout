package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/stonebirdjx/go-layout/internal/config"
	"github.com/stonebirdjx/go-layout/internal/logger"
	"github.com/stonebirdjx/go-layout/pkg/app"
)

// RunServer initializes components and starts the application lifecycle.
func RunServer(cfgFile string) error {
	// 1. & 2. Init Config
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// 3. Init Logger
	log := logger.Init(&cfg.Log)
	log.Info("Starting application bootstrap...", "env", cfg.App.Env, "version", cfg.App.Version)

	// Create App Manager
	application := app.New(
		app.WithLogger(log),
		app.WithStartTimeout(30*time.Second),
		app.WithStopTimeout(30*time.Second),
	)

	// 4. Init Clients (DB, Cache, MQ, etc.)
	application.AppendHook(app.Hook{
		Name: "database",
		OnStart: func(ctx context.Context) error {
			log.Info("Initializing database connection pool...")
			// TODO: Initialize db (e.g. gorm/sqlx), test ping
			time.Sleep(200 * time.Millisecond) // Simulated
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Closing database connections...")
			// TODO: Close db pool
			return nil
		},
	})

	application.AppendHook(app.Hook{
		Name: "redis",
		OnStart: func(ctx context.Context) error {
			log.Info("Initializing Redis client...")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Closing Redis connection...")
			return nil
		},
	})

	// 5. Init Background Tasks
	application.AppendHook(app.Hook{
		Name: "background-tasks",
		OnStart: func(ctx context.Context) error {
			log.Info("Starting background tasks scheduler...")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Stopping background tasks...")
			return nil
		},
	})

	// 6. Init Servers (HTTP / gRPC)
	httpSrv := &http.Server{
		Addr: cfg.Server.HTTP.Addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello from go-layout server!"))
		}),
	}
	application.AppendHook(app.Hook{
		Name: "http-server",
		OnStart: func(ctx context.Context) error {
			log.Info("Starting HTTP server", "addr", cfg.Server.HTTP.Addr)
			go func() {
				if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Error("HTTP server failed", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Shutting down HTTP server...")
			return httpSrv.Shutdown(ctx)
		},
	})

	// 7. Final Check & Block on Run
	log.Info("All components registered. Starting lifecycle loop...")
	
	if err := application.Run(); err != nil {
		log.Error("Application stopped with error", "error", err)
		return err
	}

	log.Info("Application stopped cleanly")
	return nil
}
