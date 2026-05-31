package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/aeschyllus/sweldo-rest/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type application struct {
	config config.Config
	db     *pgxpool.Pool
	server *http.Server
}

func (app *application) run(h http.Handler) error {
	app.server = &http.Server{
		Addr:         app.config.Addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	slog.Info("server started", "addr", app.config.Addr)

	return app.server.ListenAndServe()
}
