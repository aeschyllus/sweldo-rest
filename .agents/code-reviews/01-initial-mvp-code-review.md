# Code Review: Sweldo REST API

---

## 1. Security Issues

### DSN logged in plaintext
**File:** `cmd/main.go:33`

The full database connection string (including password) is logged on startup. An operator scrolling logs sees the DB password.

```go
logger.Info("Connected to database", "dsn", cfg.db.dsn)
```

**Fix:** Log only the host, or strip credentials from the DSN.

---

### No request body size limit
All `json.Read()` calls accept unbounded request bodies. A malicious client can send a multi-gigabyte payload and exhaust RAM.

**Fix:** Wrap the body with `http.MaxBytesReader(w, r.Body, 1<<20)` (or similar limit) before decoding.

---

### Internal error messages leaked to client
Every handler returns `err.Error()` directly to the client on 500 errors. This exposes database errors, constraint names, and internal details:

```go
json.WriteError(w, http.StatusInternalServerError, err.Error())
```

**Fix:** Log the real error server-side, return a generic "internal server error" message to the client.

---

## 2. Critical Architecture Issues

### Single DB connection — no pool
**File:** `cmd/main.go:27`

Uses `pgx.Connect()` which returns a single `*pgx.Conn`. A web server handling concurrent requests will serialize all database access. This is the biggest performance bottleneck.

```go
conn, err := pgx.Connect(ctx, cfg.db.dsn)
```

**Fix:** Use `pgxpool.New(ctx, dsn)` instead. The sqlc `Querier` interface accepts `pgxpool.Pool` as a `DBTX`.

---

### No graceful shutdown
**File:** `cmd/api.go:74`

Calls `srv.ListenAndServe()` which blocks forever. There is no `signal.Notify` catching SIGINT/SIGTERM, no `srv.Shutdown()`, no draining of in-flight requests.

**Fix:** Add a signal handler that calls `srv.Shutdown(ctx)` with a timeout.

---

### Mixed logging — `log` + `slog`
**File:** `cmd/api.go:72`

```go
log.Printf("Server has started at addr %s", app.config.addr)
```

But `main.go:23` sets up `slog` as the default logger. Two logging paths, inconsistent output format.

**Fix:** Use `slog` everywhere.

---

## 3. Bugs

### `pgconvert.FromDate` has wrong Go format string
**File:** `internal/pkg/pgconvert/date.go:20`

```go
return d.Time.Format("2026-01-02")  // should be "2006-01-02"
```

`"2026"` is not a Go reference-time token. It happens to work because Go outputs unknown tokens literally, but this is fragile and wrong.

---

### `pgconvert.ToNumeric` ignores errors
**File:** `internal/pkg/pgconvert/numeric.go:7-8`

```go
n.Scan(f)  // error silently dropped
```

A `NaN`, `Inf`, or other invalid float64 silently produces a zero/invalid `pgtype.Numeric`.

**Fix:** Return an error, or at minimum handle `NaN`/`Inf` upfront.

---

### `float64` for monetary values
Every monetary field (`BaseSalary`, `GrossPay`, `TaxDeduction`, `NetPay`, `TotalPay`) uses `float64` -> `pgtype.Numeric` conversion. Floats cannot represent exact decimal values like `0.01`. You will get rounding errors in financial calculations.

**Fix:** Store as `int64` (cents) or use a proper `decimal` type. Accept strings in JSON and parse to a `decimal.Decimal` library.

---

### GET endpoint with request body
**File:** `internal/modules/employees/handler.go:73-96`

`FindEmployeeByID` reads `company_id` from a JSON **body** on a GET request. HTTP GET requests with bodies are unusual and many proxies/balancers strip them.

```go
var req findEmployeeRequest   // has CompanyID
if err := json.Read(r, &req); err != nil { ... }
```

**Fix:** Pass `company_id` as a query parameter (like the list endpoint does).

---

### Dead struct field
**File:** `internal/modules/employees/types.go:34`

`findEmployeeRequest` has an `ID` field that is never populated (the actual ID comes from the URL param via `chi.URLParam`).

---

## 4. Code Quality

### No input validation
- `EmploymentType` / `SalaryType` are raw strings with no enum validation
- Payroll amounts can be negative
- No string length limits
- No required-field checks beyond "is it empty"

### Redundant null check
**File:** `payrolls/handler.go:59-61`

The handler calls `parseListPayrollRunsQuery` and then separately checks `query.CompanyID == 0`. The parse function already returns zero for missing values.

```go
if query.CompanyID == 0 {
    json.WriteError(w, http.StatusBadRequest, "company_id is required")
```

This logic should live inside the parse function, not duplicated in the handler.

---

### Chi RequestID middleware comment is misleading
**File:** `cmd/api.go:34`

```go
r.Use(middleware.RequestID) // Important for rate limiting
```

RequestID has nothing to do with rate limiting. This comment is wrong.

---

### `env.GetString` can't distinguish empty from unset
**File:** `internal/pkg/env/env.go:5-9`

Returns fallback for both empty string and unset variable. If a variable should be explicitly set to empty, this breaks.

---

### Healthcheck doesn't check DB
**File:** `cmd/api.go:42`

`/healthcheck` always returns "all good" even if the database is unreachable. A healthcheck should `Ping()` the DB.

---

### No pagination on employee listing
Already TODO'd in the code but not implemented. The companies endpoint supports pagination, employees don't.

---

### No CORS middleware
If consumed by a browser frontend, all requests will be blocked.

---

## 5. Missing Observability

- No structured logging with request context (request ID is set by middleware but never logged with handler output)
- No metrics (request count, latency, error rate)
- No tracing
- No `slog` attributes on server start (listen addr is logged via `log.Printf`)

---

## 6. Testing

**Zero tests.** Not a single `_test.go` file exists. The project has no:
- Unit tests for services
- Integration tests for handlers
- SQL/sqlc query tests
- Test database setup

---

## 7. Dependency Management

The `.gitignore` says it ignores "generated sqlc Go files", but those files **are** checked into git. If they're deliberately checked in (which is standard for sqlc), the gitignore rule is misleading.

---

## 8. Project Structure Suggestions

The current structure is actually quite good — clean module separation, layered architecture, sqlc for type-safe SQL. Minor suggestions:

### Suggested structure

```
cmd/
  api.go              # server setup
  main.go             # entry point
internal/
  config/
    config.go          # move config struct + parsing here
  modules/
    companies/
    employees/
    payrolls/
  adapters/
    postgresql/
      migrations/
      sqlc/
  pkg/
    env/
    json/
    pgconvert/
```

### Specific structural improvements

1. **Move `config`, `application`, `dbConfig` out of `cmd/`** — they're not command-specific. Put them in `internal/config/`.

2. **Move router `mount()` into a dedicated `internal/router/` package** — `cmd/api.go` mixes server lifecycle with route wiring.

3. **Reconsider the `pgconvert` package** — having a whole package for adapter functions suggests the domain types might not be well-separated from the DB types. Ideally your domain doesn't need to know about `pgtype.Numeric`. Consider mapping fully to domain DTOs in the service layer.

4. **The `internal/pkg/` path is slightly redundant** — `internal/` already implies private code. Could simplify to `internal/json/`, `internal/env/`, etc.

---

## Summary of Quick Wins

| Priority | Issue | File | Line |
|----------|-------|------|------|
| 🔴 Critical | Single DB connection (no pool) | `cmd/main.go` | 27 |
| 🔴 Critical | No graceful shutdown | `cmd/api.go` | 74 |
| 🔴 Critical | DSN logged with password | `cmd/main.go` | 33 |
| 🔴 Critical | GET endpoint reads request body | `employees/handler.go` | 82-85 |
| 🟡 High | `FromDate` format string bug | `pgconvert/date.go` | 20 |
| 🟡 High | `float64` for money will lose precision | `employees/types.go` | 30 |
| 🟡 High | Error messages leaked to client | all handlers | |
| 🟡 High | No request body size limit | `json/json.go` | 15 |
| 🟡 High | `ToNumeric` ignores scan errors | `pgconvert/numeric.go` | 7 |
| 🟡 High | No tests | — | |
| 🟢 Medium | Mixed `log`/`slog` | `cmd/api.go` | 72 |
| 🟢 Medium | Healthcheck doesn't ping DB | `cmd/api.go` | 41-43 |
| 🟢 Medium | No input validation | all handlers | |
| 🟢 Medium | Dead field in `findEmployeeRequest` | `employees/types.go` | 34 |
| 🔵 Low | Wrong comment on RequestID middleware | `cmd/api.go` | 34 |
| 🔵 Low | Duplicate null check | `payrolls/handler.go` | 59-61 |
