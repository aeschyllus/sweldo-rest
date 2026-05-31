package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/aeschyllus/sweldo-rest/internal/pkg/env"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	cfg := config{
		addr: ":8080",
		db: dbConfig{
			dsn: env.GetString("GOOSE_DBSTRING", "host=localhost user=postgres password=postgres dbname=sweldo sslmode=disable"),
		},
	}

	// Logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Database pool
	pool, err := pgxpool.New(ctx, cfg.db.dsn)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	logger.Info("Connected to database", "dsn", sanitizeDSN(cfg.db.dsn))

	api := application{
		config: cfg,
		db:     pool,
	}

	// Start server in background
	errCh := make(chan error, 1)
	go func() {
		if err := api.run(api.mount()); err != nil {
			errCh <- err
		}
	}()

	// Wait for interrupt signal or server error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-quit:
		logger.Info("shutting down server", "signal", sig)
	case err := <-errCh:
		logger.Error("server error", "error", err)
		os.Exit(1)
	}

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := api.server.Shutdown(shutdownCtx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
	}

	pool.Close()
	logger.Info("server stopped")
}

var dsnPasswordRe = regexp.MustCompile(`password=\S+`)

func sanitizeDSN(dsn string) string {
	return dsnPasswordRe.ReplaceAllString(dsn, "password=****")
}
