package router

import (
	"context"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/aeschyllus/sweldo-rest/internal/adapters/postgresql/sqlc"
	authmodule "github.com/aeschyllus/sweldo-rest/internal/modules/auth"
	"github.com/aeschyllus/sweldo-rest/internal/modules/companies"
	"github.com/aeschyllus/sweldo-rest/internal/modules/employees"
	"github.com/aeschyllus/sweldo-rest/internal/modules/payrolls"
	"github.com/aeschyllus/sweldo-rest/internal/pkg/auth"
	"github.com/aeschyllus/sweldo-rest/internal/pkg/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ctxKey string

const loggerCtxKey ctxKey = "logger"

func New(pool *pgxpool.Pool, jwtSecret string) http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(requestIDToSlog)
	r.Use(middleware.Logger)
	r.Use(recoveryMiddleware)
	r.Use(middleware.Timeout(time.Minute))

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		if err := pool.Ping(r.Context()); err != nil {
			slog.ErrorContext(r.Context(), "healthcheck failed", "error", err)
			json.WriteError(w, http.StatusServiceUnavailable, "database unreachable")
			return
		}
		w.Write([]byte("all good"))
	})

	// Auth module (public routes)
	authSvc := authmodule.NewService(sqlc.New(pool), jwtSecret)
	authH := authmodule.NewHandler(authSvc)
	authH.RegisterRoutes(r)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(auth.Middleware(jwtSecret))

		// Companies module
		companyService := companies.NewService(sqlc.New(pool))
		companyHandler := companies.NewHandler(companyService)
		companyHandler.RegisterRoutes(r)

		// Employees module
		employeeService := employees.NewService(sqlc.New(pool))
		employeeHandler := employees.NewHandler(employeeService)
		employeeHandler.RegisterRoutes(r)

		// Payrolls module
		payrollService := payrolls.NewService(sqlc.New(pool))
		payrollHandler := payrolls.NewHandler(payrollService)
		payrollHandler.RegisterRoutes(r)
	})

	return r
}

func requestIDToSlog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := middleware.GetReqID(r.Context())
		logger := slog.With("request_id", reqID)
		ctx := context.WithValue(r.Context(), loggerCtxKey, logger)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.ErrorContext(r.Context(), "panic recovered", "panic", rec, "stack", string(debug.Stack()))
				json.WriteError(w, http.StatusInternalServerError, "internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
