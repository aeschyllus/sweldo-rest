package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/aeschyllus/sweldo-rest/internal/adapters/postgresql/sqlc"
	"github.com/aeschyllus/sweldo-rest/internal/modules/companies"
	"github.com/aeschyllus/sweldo-rest/internal/modules/employees"
	"github.com/aeschyllus/sweldo-rest/internal/modules/payrolls"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
)

type application struct {
	config config
	db     *pgx.Conn
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID) // Injects request ID into context for logging
	r.Use(middleware.RealIP)    // Important for rate limiting and analytics
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) // Recover from crashes

	r.Use(middleware.Timeout(time.Minute))

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all good"))
	})

	// Companies module
	companyService := companies.NewService(sqlc.New(app.db))
	companyHandler := companies.NewHandler(companyService)
	companyHandler.RegisterRoutes(r)

	// Employees module
	employeeService := employees.NewService(sqlc.New(app.db))
	employeeHandler := employees.NewHandler(employeeService)
	employeeHandler.RegisterRoutes(r)

	// Payrolls module
	payrollService := payrolls.NewService(sqlc.New(app.db))
	payrollHandler := payrolls.NewHandler(payrollService)
	payrollHandler.RegisterRoutes(r)

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	slog.Info("server started", "addr", app.config.addr)

	return srv.ListenAndServe()
}
