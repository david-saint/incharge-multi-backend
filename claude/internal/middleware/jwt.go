package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/incharge/server/internal/dto"
)

type contextKey string

const (
	// UserIDKey is the context key for the authenticated user's ID.
	UserIDKey contextKey = "user_id"
	// AdminIDKey is the context key for the authenticated admin's ID.
	AdminIDKey contextKey = "admin_id"
)

// JWTAuth validates JWT tokens from the Authorization header.
func JWTAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				dto.WriteAuthError(w, "Permission Denied")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				dto.WriteAuthError(w, "Permission Denied")
				return
			}

			tokenStr := parts[1]
			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				slog.Debug("JWT validation failed", "error", err)
				dto.WriteAuthError(w, "Permission Denied")
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				dto.WriteAuthError(w, "Permission Denied")
				return
			}

			// Extract user ID from "sub" claim.
			sub, ok := claims["sub"]
			if !ok {
				dto.WriteAuthError(w, "Permission Denied")
				return
			}

			var userID uint
			switch v := sub.(type) {
			case float64:
				userID = uint(v)
			default:
				dto.WriteAuthError(w, "Permission Denied")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID extracts the authenticated user ID from context.
func GetUserID(ctx context.Context) (uint, bool) {
	id, ok := ctx.Value(UserIDKey).(uint)
	return id, ok
}

// GetAdminID extracts the authenticated admin ID from context.
func GetAdminID(ctx context.Context) (uint, bool) {
	id, ok := ctx.Value(AdminIDKey).(uint)
	return id, ok
}
