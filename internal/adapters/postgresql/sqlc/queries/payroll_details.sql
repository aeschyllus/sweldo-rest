-- name: CreatePayrollDetail :one
INSERT INTO payroll_details (payroll_run_id, employee_id, gross_pay, tax_deduction, net_pay, hourly_rate, hours_worked, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: FindPayrollDetailByID :one
SELECT * FROM payroll_details WHERE id = $1;

-- name: ListAllPayrollDetailsByRunID :many
SELECT * FROM payroll_details WHERE payroll_run_id = $1;

-- name: ListAllPayrollDetailsByEmployeeID :many
SELECT *
FROM payroll_details
WHERE employee_id = $1
ORDER BY created_at DESC;
