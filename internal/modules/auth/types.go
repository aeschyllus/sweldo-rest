package auth

import (
	"context"

	"github.com/aeschyllus/sweldo-rest/internal/adapters/postgresql/sqlc"
)

type handler struct {
	service Service
}

type service struct {
	repo      sqlc.Querier
	jwtSecret string
}

type Service interface {
	Register(ctx context.Context, params RegisterParams) (AuthResponse, error)
	Login(ctx context.Context, params LoginParams) (AuthResponse, error)
}

type RegisterParams struct {
	CompanyID int64
	Email     string
	Password  string
	CreatedBy int64
}

type LoginParams struct {
	Email    string
	Password string
}

type AuthResponse struct {
	Token     string `json:"token"`
	UserID    int64  `json:"user_id"`
	CompanyID int64  `json:"company_id"`
	Email     string `json:"email"`
}

type registerRequest struct {
	CompanyID int64  `json:"company_id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
