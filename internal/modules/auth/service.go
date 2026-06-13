package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/aeschyllus/sweldo-rest/internal/adapters/postgresql/sqlc"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func NewService(repo sqlc.Querier, jwtSecret string) Service {
	return &service{repo: repo, jwtSecret: jwtSecret}
}

func (s *service) Register(ctx context.Context, params RegisterParams) (AuthResponse, error) {
	if params.Email == "" {
		return AuthResponse{}, fmt.Errorf("email is required")
	}
	if params.Password == "" {
		return AuthResponse{}, fmt.Errorf("password is required")
	}
	if params.CompanyID == 0 {
		return AuthResponse{}, fmt.Errorf("company_id is required")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.repo.CreateUser(ctx, sqlc.CreateUserParams{
		CompanyID:    params.CompanyID,
		Email:        params.Email,
		PasswordHash: string(hash),
	})
	if err != nil {
		return AuthResponse{}, fmt.Errorf("failed to create user: %w", err)
	}

	token, err := s.generateToken(user)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("failed to generate token: %w", err)
	}

	return AuthResponse{
		Token:     token,
		UserID:    user.ID,
		CompanyID: user.CompanyID,
		Email:     user.Email,
	}, nil
}

func (s *service) Login(ctx context.Context, params LoginParams) (AuthResponse, error) {
	if params.Email == "" {
		return AuthResponse{}, fmt.Errorf("email is required")
	}
	if params.Password == "" {
		return AuthResponse{}, fmt.Errorf("password is required")
	}

	user, err := s.repo.FindUserByEmail(ctx, params.Email)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(params.Password)); err != nil {
		return AuthResponse{}, fmt.Errorf("invalid email or password")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("failed to generate token: %w", err)
	}

	return AuthResponse{
		Token:     token,
		UserID:    user.ID,
		CompanyID: user.CompanyID,
		Email:     user.Email,
	}, nil
}

func (s *service) generateToken(user sqlc.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    user.ID,
		"company_id": user.CompanyID,
		"email":      user.Email,
		"exp":        time.Now().Add(72 * time.Hour).Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
