-- name: CreatePayrollRun :one
INSERT INTO payroll_runs (company_id, run_date, total_employees, total_pay, created_by)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListPayrollRunsByCompanyID :many
SELECT *
FROM payroll_runs
WHERE company_id = $1
ORDER BY run_date DESC;

-- name: FindPayrollRunByID :one
SELECT * FROM payroll_runs WHERE id = $1;

-- name: UpdatePayrollRunByID :one
UPDATE payroll_runs
SET total_employees = $1, total_pay = $2, updated_by = $3, updated_at = NOW()
WHERE id = $4
RETURNING *;

-- name: FinalizePayrollRunByID :one
UPDATE payroll_runs
SET status = 'FINALIZED', updated_by = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;
