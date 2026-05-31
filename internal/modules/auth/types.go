package auth

import (
	"context"
)

type handler struct {
	service Service
}

type service struct {
	repo      UserRepository
	jwtSecret string
}

type Service interface {
	Register(ctx context.Context, params RegisterParams) (AuthResponse, error)
	Login(ctx context.Context, params LoginParams) (AuthResponse, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, params CreateUserParams) (User, error)
	FindUserByEmail(ctx context.Context, email string) (*User, error)
}

type User struct {
	ID           int64
	CompanyID    int64
	Email        string
	PasswordHash string
}

type CreateUserParams struct {
	CompanyID    int64
	Email        string
	PasswordHash string
}

type RegisterParams struct {
	CompanyID int64
	Email     string
	Password  string
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
