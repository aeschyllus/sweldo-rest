-- name: CreateUser :one
INSERT INTO users (company_id, email, password_hash, created_by)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: FindUserByEmail :one
SELECT * FROM users WHERE email = $1;
