package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"

	"github.com/aeschyllus/sweldo-rest/internal/config"
	"github.com/aeschyllus/sweldo-rest/internal/pkg/env"
	"github.com/aeschyllus/sweldo-rest/internal/router"
)

func main() {
	ctx := context.Background()

	env.LoadEnvFile(".env")

	cfg := config.Config{
		Addr: env.GetString("PORT", ":8080"),
		DB: config.DBConfig{
			DSN: env.GetString("DB_DSN", env.GetString("GOOSE_DBSTRING", "host=localhost user=postgres password=postgres dbname=sweldo sslmode=disable")),
		},
	}

	jwtSecret := env.GetString("JWT_SECRET", "dev-secret-change-in-production")

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	pool, err := pgxpool.New(ctx, cfg.DB.DSN)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	logger.Info("Connected to database", "dsn", sanitizeDSN(cfg.DB.DSN))

	sqlDB, err := sql.Open("pgx", cfg.DB.DSN)
	if err != nil {
		panic(err)
	}

	if err := goose.Up(sqlDB, "internal/adapters/postgresql/migrations"); err != nil {
		logger.Error("migration failed", "error", err)
		panic(err)
	}
	sqlDB.Close()

	logger.Info("Migrations applied successfully")

	api := application{
		config: cfg,
		db:     pool,
	}

	errCh := make(chan error, 1)
	go func() {
		if err := api.run(router.New(pool, jwtSecret)); err != nil {
			errCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-quit:
		logger.Info("shutting down server", "signal", sig)
	case err := <-errCh:
		logger.Error("server error", "error", err)
		os.Exit(1)
	}

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
