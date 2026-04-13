package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateTestToken(secret string, userID uint, ttl time.Duration) string {
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(ttl).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, _ := token.SignedString([]byte(secret))
	return str
}

func TestJWTAuthMiddleware(t *testing.T) {
	secret := "jwt-test-secret"

	t.Run("valid token sets user ID in context", func(t *testing.T) {
		token := generateTestToken(secret, 42, 60*time.Minute)

		var capturedID uint
		var capturedOK bool
		handler := JWTAuth(secret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedID, capturedOK = GetUserID(r.Context())
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
		if !capturedOK || capturedID != 42 {
			t.Fatalf("expected userID=42, got %d (ok=%v)", capturedID, capturedOK)
		}
	})

	t.Run("missing Authorization header returns 401", func(t *testing.T) {
		handler := JWTAuth(secret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", rec.Code)
		}
	})

	t.Run("malformed Authorization header returns 401", func(t *testing.T) {
		handler := JWTAuth(secret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "NotBearer token")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		// "NotBearer" is not "Bearer", should fail
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", rec.Code)
		}
	})

	t.Run("expired token returns 401", func(t *testing.T) {
		token := generateTestToken(secret, 42, -1*time.Minute)

		handler := JWTAuth(secret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", rec.Code)
		}
	})

	t.Run("wrong secret returns 401", func(t *testing.T) {
		token := generateTestToken("other-secret", 42, 60*time.Minute)

		handler := JWTAuth(secret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", rec.Code)
		}
	})
}

func TestGetUserID(t *testing.T) {
	t.Run("returns user ID from context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), UserIDKey, uint(99))
		id, ok := GetUserID(ctx)
		if !ok || id != 99 {
			t.Fatalf("expected 99, got %d (ok=%v)", id, ok)
		}
	})

	t.Run("returns false when not set", func(t *testing.T) {
		_, ok := GetUserID(context.Background())
		if ok {
			t.Fatal("expected ok=false when user ID not in context")
		}
	})
}

func TestGetAdminID(t *testing.T) {
	t.Run("returns admin ID from context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), AdminIDKey, uint(5))
		id, ok := GetAdminID(ctx)
		if !ok || id != 5 {
			t.Fatalf("expected 5, got %d (ok=%v)", id, ok)
		}
	})

	t.Run("returns false when not set", func(t *testing.T) {
		_, ok := GetAdminID(context.Background())
		if ok {
			t.Fatal("expected ok=false when admin ID not in context")
		}
	})
}
