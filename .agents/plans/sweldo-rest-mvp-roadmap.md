# sweldo-rest — MVP Roadmap

> A fully functional, deployable payroll REST API covering companies, employees, and payroll runs/details — with authentication in place and all known bugs resolved.

---

## Phase 1 — Bug Fixes & Cleanup

> Resolve existing code debt before building further. These are blocking issues surfaced in the Apr 4 code review.

### companies module
- Fix the missing `return` after the error branch in `ListCompanies` handler (double-write bug)

### employees module
- Fix `FindEmployeeByID` — remove the `json.Read` body call, parse ID from URL param only
- Remove the dead `findEmployeeRequest.ID` field
- Fix `UpdateEmployeeByID` to scope by `company_id` (even if hardcoded for now, until JWT lands)
- Standardize `company_id` sourcing — consistently use query param across all read endpoints

---

## Phase 2 — Authentication & Identity

> This unlocks the `company_id` TODO that has been deferred across the entire codebase.

- Add a JWT middleware (e.g. `golang-jwt/jwt`)
- Extract `company_id` (and eventually `user_id`) from JWT claims in middleware
- Thread `company_id` into the request context so handlers can pull it without touching the body or query params
- Remove all manual `company_id` passing from request bodies/query params on write endpoints
- Scope all write operations (`CreateEmployee`, `UpdateEmployeeByID`, etc.) by the JWT-extracted `company_id`

**Suggested location:** `internal/pkg/auth/` or `internal/pkg/middleware/`

---

## Phase 3 — Payrolls Module Finalization

> The module was generated but a few loose ends remain to be confirmed and implemented.

- Confirm whether the unique constraint on `payroll_details(payroll_run_id, employee_id)` was adopted in the migration
- Confirm whether the `status` field and `employee_id` index were added
- Add a `PATCH /payroll-runs/{runID}/finalize` endpoint — once a run is finalized, details should be immutable; enforce this with a status check in the service layer before any detail write
- Wire up the `payrolls` module routes into `cmd/api.go`

---

## Phase 4 — Payroll Business Logic

> The core domain work that makes this useful as a payroll API.

- **Gross pay calculation** — derive gross from rate × hours (or salary-based) when creating a payroll detail
- **Deductions** — add a `deductions` table or embed deduction fields in `payroll_details` (SSS, PhilHealth, Pag-IBIG for Philippine payroll)
- **Net pay** — calculated field: `gross_pay - total_deductions`
- **Package decision** — if deductions have no standalone lifecycle, keep them inside the `payrolls` package (apply the same consolidation rule used for `payroll_details`)

---

## Phase 5 — Observability & Hardening

> Before calling it shippable.

- **Structured logging** — log request method, path, status code, and latency (e.g. `slog` or `zerolog`)
- **Centralized error handling** — middleware that produces a consistent error response shape across all endpoints
- **Request validation** — required fields, value ranges on all request types
- **Health check** — `GET /health` that pings the DB and returns status
- **Environment config** — ensure `internal/pkg/env` covers all required vars (DB DSN, JWT secret, port)

---

## Phase 6 — Deployment

- Multi-stage `Dockerfile` (build → minimal runtime image)
- `docker-compose.yml` for local dev (API + PostgreSQL)
- CI pipeline — at minimum: `go build`, `go vet`, `go test ./...`
- Migration execution on startup or as a separate init step via goose

---

## Summary

| Phase | Focus | Depends On |
|---|---|---|
| 1 | Bug fixes in companies/employees | — |
| 2 | JWT auth + `company_id` from context | Phase 1 |
| 3 | Payrolls finalization + immutability | Phase 2 |
| 4 | Payroll business logic (gross/net/deductions) | Phase 3 |
| 5 | Logging, validation, error handling | Phase 3 |
| 6 | Docker, CI, deployment | Phases 4–5 |

> **Highest-leverage item:** Phase 2 (JWT auth) unblocks every scoping and security concern that has been deferred with TODOs across the codebase.

---

## Implementation Status (Assessed 2026-05-09)

> Each item is marked ✅ (done), ❌ (not done), or ⚠️ (partial).

### Phase 1 — Bug Fixes & Cleanup — ⚠️ PARTIALLY DONE

**companies module**
- ✅ Fix missing `return` after error branch in `ListCompanies` — handler at `companies/handler.go:51-66` properly returns after both error branches.

**employees module**
- ❌ `FindEmployeeByID` still calls `json.Read` on the request body at `employees/handler.go:84` to read `findEmployeeRequest`. The phase requirement is to remove the body read entirely and parse only the URL param.
- ✅ Dead `findEmployeeRequest.ID` field was removed — struct at `employees/types.go:33-35` only has `CompanyID`.
- ❌ `UpdateEmployeeByID` is not scoped by `company_id` — `UpdateEmployeeParams` at `employees/types.go:65-72` and the sqlc call at `employees/service.go:50-57` both lack `CompanyID`.
- ❌ `company_id` sourcing is inconsistent — `ListEmployeesByCompanyID` uses a query param, but `CreateEmployee` and `FindEmployeeByID` read it from the request body.

### Phase 2 — Authentication & Identity — ❌ NOT STARTED

- ❌ No JWT library in `go.mod` — `github.com/golang-jwt/jwt` is absent.
- ❌ No auth middleware exists anywhere (`internal/pkg/auth/` or `internal/pkg/middleware/` don't exist).
- ❌ No login/register endpoint.
- ❌ `company_id` is still passed manually via request body and query params throughout all handlers — no JWT extraction, no context threading.
- ✅ Several TODO comments acknowledge this is deferred (e.g. `employees/handler.go:52`).

### Phase 3 — Payrolls Module Finalization — ⚠️ PARTIALLY DONE

- ✅ Unique constraint `idx_payroll_details_run_employee` on `(payroll_run_id, employee_id)` is present in migration `00001_init.sql:58`.
- ⚠️ `employee_id` index `idx_payroll_details_employee_id` is present (`00001_init.sql:56`). ❌ `status` field does **not** exist on `payroll_runs` or `payroll_details` — neither in the migration schema nor in any Go type.
- ❌ No `PATCH /payroll-runs/{runID}/finalize` endpoint exists — `payrolls/handler.go:16-30` has no such route.
- ✅ Payroll routes are wired into `cmd/api.go:56-58`.

### Phase 4 — Payroll Business Logic — ❌ NOT STARTED

- ❌ `gross_pay` is a user-supplied input field (`createPayrollDetailRequest.GrossPay`) — no computation from rate × hours or salary.
- ❌ No deductions table or embedded deduction fields exist beyond a single `tax_deduction` column. No SSS, PhilHealth, or Pag-IBIG support.
- ❌ `net_pay` is a user-supplied input — not derived from `gross_pay - total_deductions`.
- ❌ Package decision (deductions placement) is moot — no deduction logic exists yet.

### Phase 5 — Observability & Hardening — ⚠️ PARTIALLY DONE

- ✅ Structured logging is in place — `slog` configured in `cmd/main.go:24-25`, chi `Logger` middleware active in `cmd/api.go:36`, `slog.ErrorContext` used in handlers.
- ❌ No centralized error handling middleware. Each handler calls `json.WriteError` individually. There's no consistent error response envelope.
- ❌ No request validation beyond JSON parse failures. No required-field checks, no value-range enforcement.
- ⚠️ `GET /healthcheck` exists at `cmd/api.go:41-43` but returns a static string — it does **not** ping the database or report actual status.
- ⚠️ `internal/pkg/env/env.go` only exposes `GetString(string, string) string`. No `GetInt`, `GetBool`, no validation, no `.env` file loading. Only one env var (`GOOSE_DBSTRING`) is read.

### Phase 6 — Deployment — ❌ NOT STARTED

- ❌ No `Dockerfile` exists anywhere in the repository.
- ⚠️ `docker-compose.yaml` exists but only defines a Postgres service — no app service is included.
- ❌ No `.github/workflows/` directory — no CI pipeline of any kind. (Agent skill templates exist under `.agents/skills/golang-continuous-integration/assets/` but are not active.)
- ❌ Goose migration execution is not integrated into `main.go` — migrations must be run manually via the `goose` CLI.

### Summary Table

| Phase | Focus | Status |
|---|---|---|
| 1 | Bug fixes in companies/employees | ⚠️ Partial — companies fix done, most employees fixes not |
| 2 | JWT auth + `company_id` from context | ❌ Not started |
| 3 | Payrolls finalization + immutability | ⚠️ Partial — routes wired, constraints present; no `status` field, no finalize endpoint |
| 4 | Payroll business logic (gross/net/deductions) | ❌ Not started |
| 5 | Logging, validation, error handling | ⚠️ Partial — slog + health endpoint exist; no centralized errors, validation, or env config |
| 6 | Docker, CI, deployment | ❌ Not started — no Dockerfile, no CI, no migration startup |
