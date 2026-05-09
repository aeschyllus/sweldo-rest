-- name: CreatePayrollDetail :one
INSERT INTO payroll_details (payroll_run_id, employee_id, gross_pay, tax_deduction, net_pay)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListAllPayrollDetailsByRunID :many
SELECT * FROM payroll_details WHERE payroll_run_id = $1;

-- name: ListAllPayrollDetailsByEmployeeID :many
SELECT *
FROM payroll_details
WHERE employee_id = $1
ORDER BY created_at DESC;