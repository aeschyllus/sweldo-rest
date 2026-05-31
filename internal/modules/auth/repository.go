package auth

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgUserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &pgUserRepo{pool: pool}
}

const createUserSQL = `
INSERT INTO users (company_id, email, password_hash)
VALUES ($1, $2, $3)
RETURNING id, company_id, email, password_hash
`

func (r *pgUserRepo) CreateUser(ctx context.Context, params CreateUserParams) (User, error) {
	row := r.pool.QueryRow(ctx, createUserSQL, params.CompanyID, params.Email, params.PasswordHash)
	var u User
	err := row.Scan(&u.ID, &u.CompanyID, &u.Email, &u.PasswordHash)
	if err != nil {
		return User{}, fmt.Errorf("create user: %w", err)
	}
	return u, nil
}

const findUserByEmailSQL = `
SELECT id, company_id, email, password_hash FROM users WHERE email = $1
`

func (r *pgUserRepo) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	row := r.pool.QueryRow(ctx, findUserByEmailSQL, email)
	var u User
	err := row.Scan(&u.ID, &u.CompanyID, &u.Email, &u.PasswordHash)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("find user: %w", err)
	}
	return &u, nil
}
