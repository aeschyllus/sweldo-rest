# Code Review Implementation Plan — Sweldo REST API

## Overview

Systematic remediation of issues identified in `.agents/code-reviews/01-initial-mvp-code-review.md`, organized into 6 phases with clear dependency ordering. Each phase is self-contained and can be implemented and verified independently.

---

## Phase 1 — Quick fixes (no downstream impact) ✅ DONE

These are all isolated changes that don't affect other code. Runnable in parallel.

### 1.1 DSN logged in plaintext
- **File:** `cmd/main.go:33`
- **Change:** Strip credentials from DSN before logging
- **Detail:** Parse DSN as URL, redact password portion, log sanitized version only

### 1.2 Internal error messages leaked to client
- **Files:**
  - `internal/modules/companies/handler.go` — lines 36, 58, 75, 103
  - `internal/modules/employees/handler.go` — lines 27, 40, 59, 65, 77, 83, 92, 103, 108, 122
  - `internal/modules/payrolls/handler.go` — lines 34, 45, 55, 60, 66, 77, 83, 94, 100, 109, 120, 127, 137, 149, 154, 164, 169, 174
- **Change:** Log real error server-side via `slog.ErrorContext(r.Context(), ...)`, return generic `"internal server error"` to client
- **Detail:** Replace `json.WriteError(w, http.StatusInternalServerError, err.Error())` with `json.WriteError(w, http.StatusInternalServerError, "internal server error")` on all 500 paths. Non-500 errors keep their messages.

### 1.3 No request body size limit
- **File:** `internal/pkg/json/json.go:14-18`
- **Change:** Wrap `r.Body` with `http.MaxBytesReader`
- **Detail:** `r.Body = http.MaxBytesReader(w, r.Body, 1<<20)` before decoding. Returns 413 on oversized payloads.

### 1.4 `FromDate` format string bug
- **File:** `internal/pkg/pgconvert/date.go:20`
- **Change:** `"2026-01-02"` → `"2006-01-02"`
- **Detail:** Go reference time is `Mon Jan 2 15:04:05 MST 2006`; year token must be `2006`.

### 1.5 `ToNumeric` ignores errors
- **File:** `internal/pkg/pgconvert/numeric.go:5-9`
- **Change:** Return error from `ToNumeric`, guard `NaN`/`Inf`
- **Detail:** Check `math.IsNaN(f)` and `math.IsInf(f, 0)` before `Scan`. Return `(pgtype.Numeric, error)`.

### 1.6 Mixed `log`/`slog`
- **File:** `cmd/api.go:72`
- **Change:** `log.Printf` → `slog.Info`
- **Detail:** `slog.Info("server started", "addr", app.config.addr)`

### 1.7 Wrong comment on RequestID middleware
- **File:** `cmd/api.go:34`
- **Change:** Fix misleading comment
- **Detail:** `r.Use(middleware.RequestID) // Injects request ID into context for logging`

### 1.8 Dead struct field + duplicate null check
- **Files:**
  - `internal/modules/employees/types.go:33-36` — remove `ID` field from `findEmployeeRequest` (unused; actual ID comes from URL param)
  - `internal/modules/payrolls/handler.go:59-61` — consolidate `CompanyID == 0` check into `parseListPayrollRunsQuery` (return error from parse function instead of checking after)
  - `internal/modules/payrolls/handler.go:168-170` — same pattern for EmployeeID in `ListPayrollDetailsByEmployeeID`

---

## Phase 2 — Restructure cmd/ + architecture foundation

### 2.1 Move config out of cmd/
- **New file:** `internal/config/config.go`
- **Change:** Move `config`, `dbConfig` structs from `cmd/api.go` here
- **Detail:** `cmd/` should only have entry point and server lifecycle. Config structs are application-wide.

### 2.2 Move router out of cmd/
- **New file:** `internal/router/router.go`
- **Change:** Move `mount()` from `cmd/api.go` into `internal/router/router.go`
- **Detail:** Accept config and pool as constructor params. Return `http.Handler`.

### 2.3 Single DB connection → connection pool
- **Files:** `internal/config/config.go`, `internal/router/router.go`, `cmd/main.go`
- **Change:** `pgx.Connect(ctx, dsn)` → `pgxpool.New(ctx, dsn)`; change `*pgx.Conn` → `*pgxpool.Pool`
- **Detail:** `sqlc.New()` accepts any `DBTX` interface — `*pgxpool.Pool` satisfies it. Concurrency bottleneck eliminated.

### 2.4 Graceful shutdown
- **File:** `cmd/main.go`
- **Change:** Add `signal.Notify` for `SIGINT`/`SIGTERM`, call `srv.Shutdown(ctx)` with 30s timeout
- **Detail:** Listen in a goroutine, block on signal channel, then `srv.Shutdown(ctx)` + `pool.Close()`

---

## Phase 3 — Module-level type changes (batched, single pass)

### 3.1 Remove `pgconvert` package
- **Delete:** `internal/pkg/pgconvert/` (4 files: `date.go`, `numeric.go`, `text.go`, `time.go`)
- **Change:** Inline conversions where used
- **Detail:**
  - `ToNumeric`/`FromNumeric` → replaced by `decimal.Decimal` ↔ `pgtype.Numeric` helpers (place in service files or `internal/decimal/` helper if shared)
  - `ToDate`/`FromDate` → inline in payrolls service (2 lines each)
  - `ToText` → inline in companies service (3 lines)
  - `FromTimestamptz` → inline where used (`.Time` accessor)

### 3.2 Typed enums for EmploymentType/SalaryType
- **File:** `internal/modules/employees/types.go`
- **Change:** Add string enum types with const blocks + validation
- **Detail:**
  ```go
  type EmploymentType string
  const (
      EmploymentTypeFullTime   EmploymentType = "FULL_TIME"
      EmploymentTypePartTime   EmploymentType = "PART_TIME"
      EmploymentTypeContractor EmploymentType = "CONTRACTOR"
  )
  func (e EmploymentType) Valid() bool { ... }

  type SalaryType string
  const (
      SalaryTypeMonthly      SalaryType = "MONTHLY"
      SalaryTypeHourly        SalaryType = "HOURLY"
      SalaryTypeProjectBased SalaryType = "PROJECT_BASED"
  )
  func (s SalaryType) Valid() bool { ... }
  ```

### 3.3 `decimal.Decimal` for monetary values
- **Files:**
  - `internal/modules/employees/types.go` — `BaseSalary float64` → `decimal.Decimal` in request, params, response structs
  - `internal/modules/payrolls/types.go` — `TotalPay`, `GrossPay`, `TaxDeduction`, `NetPay` same change
- **Dependency add:** `github.com/shopspring/decimal`
- **Detail:** JSON accepts string `"450.00"` via custom unmarshal or `json:"...,string"` tag. Service layer converts `decimal.Decimal` → `pgtype.Numeric` and back. No precision loss.

### 3.4 Input validation in handlers
- **Files:** All 3 handler files
- **Change:** Add validation calls after `json.Read()`
- **Detail:**
  - Call `EmploymentType(req.EmploymentType).Valid()` and return 400 if invalid
  - Call `SalaryType(req.SalaryType).Valid()` if present
  - Check required string fields non-empty with clear error message
  - All validation errors use `json.WriteError(w, http.StatusBadRequest, "descriptive message")`
  - 500 errors remain generic per 1.2

### 3.5 Pagination on employee listing
- **File:** `internal/adapters/postgresql/sqlc/queries/employees.sql`
- **Change:** Add `LIMIT`/`OFFSET` parameters matching `companies.sql` pattern
- **Detail:**
  ```sql
  ORDER BY id ASC
  LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);
  ```
- **Requires:** `sqlc generate` to regenerate `employees.sql.go` and `querier.go`
- **Files:**
  - `internal/modules/employees/query.go` — parse `limit`/`offset` from query params
  - `internal/modules/employees/types.go` — add `PageLimit`, `PageOffset` to `ListEmployeesParams`
  - `internal/modules/employees/service.go` — pass pagination through to sqlc

---

## Phase 4 — CORS + healthcheck

### 4.1 CORS middleware (allow all origins)
- **File:** `internal/router/router.go`
- **Change:** Add chi CORS middleware
- **Detail:**
  ```go
  r.Use(cors.Handler(cors.Options{
      AllowedOrigins:   []string{"*"},
      AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
      AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
      ExposedHeaders:   []string{"Link"},
      AllowCredentials: false,
      MaxAge:           300,
  }))
  ```

### 4.2 Healthcheck pings DB
- **File:** `internal/router/router.go`
- **Change:** Replace inline `"all good"` with `pool.Ping(ctx)`
- **Detail:** `err := pool.Ping(r.Context())` — return 503 if DB unreachable, 200 otherwise.

---

## Phase 5 — Observability

### 5.1 Structured logging with request context
- **File:** `internal/router/router.go`
- **Change:** Add middleware that injects request ID into slog context
- **Detail:**
  ```go
  r.Use(func(next http.Handler) http.Handler {
      return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
          reqID := middleware.GetReqID(r.Context())
          r = r.WithContext(slog.NewContext(r.Context(), slog.With("request_id", reqID)))
          next.ServeHTTP(w, r)
      })
  })
  ```
- **Files:** All 3 handler files — replace `slog.Error` with `slog.ErrorContext(r.Context(), ...)` in error paths.

---

## Phase 6 — Testing

### 6.1 Unit tests for services
- **Files:**
  - `internal/modules/companies/service_test.go`
  - `internal/modules/employees/service_test.go`
  - `internal/modules/payrolls/service_test.go`
- **How:** Mock `sqlc.Querier` interface using testify, test service methods in isolation
- **Covers:** CRUD operations, validation rules, decimal conversions, enum validation

### 6.2 Integration tests for handlers
- **Files:**
  - `internal/modules/companies/handler_test.go`
  - `internal/modules/employees/handler_test.go`
  - `internal/modules/payrolls/handler_test.go`
- **How:** Spin up test DB via `docker-compose`, run migrations, test full HTTP endpoints
- **Covers:** Request/response cycles, auth, validation error format, pagination, healthcheck

---

## Dependency Graph

```
Phase 1 (quick fixes)
  ├── independent
  ├── no changes needed after
  └── safe to run anytime

Phase 2 (restructure cmd/ + architectural)
  ├── must run BEFORE Phase 3 (both touch cmd/ structure)
  ├── modifies: cmd/main.go, cmd/api.go
  └── creates: internal/config/config.go, internal/router/router.go

Phase 3 (module-level changes)
  ├── must run AFTER Phase 2 (stable cmd/ structure)
  ├── contains sub-dependency: 3.5 (pagination sqlc) needs sqlc generate
  ├── modifies: module types.go, handler.go, service.go, query.go
  ├── deletes: internal/pkg/pgconvert/
  └── File changes ranked by risk:
      [low]   3.2 (typed enums) — additive, no behavioral change to existing code
      [medium] 3.4 (input validation) — new checks, could reject previously-accepted data
      [medium] 3.1 (pgconvert removal) — mechanical inline, testable
      [high]   3.3 (decimal.Decimal) — API contract changes (float→string in JSON)
      [high]   3.5 (pagination) — modifies sql + sqlc generation

Phase 4 (CORS + healthcheck)
  ├── must run AFTER Phase 2.2 (router exists as package)
  └── modifies: internal/router/router.go

Phase 5 (observability)
  ├── must run AFTER Phase 2.2 (router exists)
  └── modifies: internal/router/router.go, all handler files

Phase 6 (testing)
  ├── must run AFTER Phase 3 (stable module interfaces)
  └── creates: 6 test files
```

---

## File Change Summary

| Phase | New Files | Modified Files | Deleted Files |
|-------|-----------|----------------|---------------|
| 1 | — | `cmd/main.go`, `cmd/api.go`, `internal/pkg/json/json.go`, `internal/pkg/pgconvert/date.go`, `internal/pkg/pgconvert/numeric.go`, `internal/modules/employees/types.go`, `internal/modules/payrolls/handler.go` | — |
| 2 | `internal/config/config.go`, `internal/router/router.go` | `cmd/main.go` (rewrite), `cmd/api.go` (reduced) | — |
| 3 | — | `employees/types.go`, `payrolls/types.go`, `employees/handler.go`, `employees/service.go`, `payrolls/handler.go`, `payrolls/service.go`, `companies/service.go`, `employees/query.go`, `queries/employees.sql` | `internal/pkg/pgconvert/*` (4 files) |
| 4 | — | `internal/router/router.go` | — |
| 5 | — | `internal/router/router.go`, `companies/handler.go`, `employees/handler.go`, `payrolls/handler.go` | — |
| 6 | 6 test files | — | — |

**Totals:** ~13 new files, ~25 modified files, ~4 deleted files

---

## Dependencies to Add

- `github.com/shopspring/decimal` (Phase 3.3)
- `github.com/go-chi/cors` (Phase 4.1) — already part of chi ecosystem
- `github.com/stretchr/testify` (Phase 6) — dev dependency
