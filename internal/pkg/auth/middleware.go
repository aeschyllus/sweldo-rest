package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const (
	claimsCtxKey ctxKey = "auth_claims"
)

type Claims struct {
	UserID    int64  `json:"user_id"`
	CompanyID int64  `json:"company_id"`
	Email     string `json:"email"`
	jwt.RegisteredClaims
}

func Middleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenStr == authHeader {
				http.Error(w, `{"error":"invalid authorization format"}`, http.StatusUnauthorized)
				return
			}

			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), claimsCtxKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ClaimsFromContext(ctx context.Context) *Claims {
	claims, _ := ctx.Value(claimsCtxKey).(*Claims)
	return claims
}

func CompanyIDFromContext(ctx context.Context) int64 {
	if claims := ClaimsFromContext(ctx); claims != nil {
		return claims.CompanyID
	}
	return 0
}

func UserIDFromContext(ctx context.Context) int64 {
	if claims := ClaimsFromContext(ctx); claims != nil {
		return claims.UserID
	}
	return 0
}
