-- name: CreateDeduction :one
INSERT INTO deductions (payroll_detail_id, deduction_type, amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListDeductionsByPayrollDetailID :many
SELECT * FROM deductions WHERE payroll_detail_id = $1;

-- name: DeleteDeduction :one
DELETE FROM deductions WHERE id = $1 RETURNING *;
