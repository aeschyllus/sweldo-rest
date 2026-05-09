package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/aeschyllus/sweldo-rest/internal/pkg/env"
	"github.com/jackc/pgx/v5"
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

	// Database
	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx) // Close DB connection when app shuts down

	logger.Info("Connected to database", "dsn", cfg.db.dsn)

	api := application{
		config: cfg,
		db:     conn,
	}

	if err := api.run(api.mount()); err != nil {
		logger.Error("Server has failed to start", "error", err)
		os.Exit(1)
	}
}
