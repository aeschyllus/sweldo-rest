-- name: CreateEmployee :one
INSERT INTO employees (company_id, first_name, last_name, employment_type, salary_type, base_salary)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListEmployeesByCompanyID :many
SELECT *
FROM employees
WHERE company_id = $1
AND (
    sqlc.narg(name)::text IS NULL OR
    first_name ILIKE '%' || sqlc.narg(name) || '%' OR
    last_name ILIKE '%' || sqlc.narg(name) || '%'
)
ORDER BY id ASC
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: FindEmployeeByID :one
SELECT * FROM employees WHERE id = $1 AND company_id = $2;

-- name: UpdateEmployeeByID :one
UPDATE employees
SET first_name = $1, last_name = $2, employment_type = $3, salary_type = $4, base_salary = $5, updated_at = NOW()
WHERE id = $6 AND company_id = $7
RETURNING *;