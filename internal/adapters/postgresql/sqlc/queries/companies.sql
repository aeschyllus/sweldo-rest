-- name: CreateCompany :one
INSERT INTO companies (name, tax_id)
VALUES ($1, $2)
RETURNING *;

-- name: ListCompanies :many
SELECT *
FROM companies
WHERE (sqlc.narg(name)::text IS NULL OR name ILIKE '%' || sqlc.narg(name) || '%')
ORDER BY id
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: FindCompanyByID :one
SELECT * FROM companies WHERE id = $1;

-- name: UpdateCompanyByID :one
UPDATE companies
SET name = $1, tax_id = $2, updated_at = NOW()
WHERE id = $3
RETURNING *;