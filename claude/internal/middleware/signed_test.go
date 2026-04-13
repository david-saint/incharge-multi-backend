package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGenerateSignedURL(t *testing.T) {
	secret := "test-secret-key"

	t.Run("produces valid URL format", func(t *testing.T) {
		url := GenerateSignedURL("http://localhost:8080/api/v1/user/email/verify", 42, secret, 60)
		if !strings.Contains(url, "/42?") {
			t.Fatalf("URL should contain user ID: %s", url)
		}
		if !strings.Contains(url, "expires=") {
			t.Fatalf("URL should contain expires param: %s", url)
		}
		if !strings.Contains(url, "signature=") {
			t.Fatalf("URL should contain signature param: %s", url)
		}
	})

	t.Run("different users produce different URLs", func(t *testing.T) {
		url1 := GenerateSignedURL("http://localhost/verify", 1, secret, 60)
		url2 := GenerateSignedURL("http://localhost/verify", 2, secret, 60)
		if url1 == url2 {
			t.Fatal("different users should produce different URLs")
		}
	})
}

func TestSignedURLMiddleware(t *testing.T) {
	secret := "test-secret-key"

	t.Run("valid signature passes through", func(t *testing.T) {
		baseURL := "http://localhost:8080/api/v1/user/email/verify"
		signedURL := GenerateSignedURL(baseURL, 42, secret, 60)

		// Parse the generated URL to get the path and query.
		req := httptest.NewRequest(http.MethodGet, signedURL, nil)
		rec := httptest.NewRecorder()

		handler := SignedURL(secret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d, body: %s", rec.Code, rec.Body.String())
		}
	})

	t.Run("missing signature returns 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/user/email/verify/42?expires=9999999999", nil)
		rec := httptest.NewRecorder()

		handler := SignedURL(secret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", rec.Code)
		}
	})

	t.Run("missing expires returns 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/user/email/verify/42?signature=abc", nil)
		rec := httptest.NewRecorder()

		handler := SignedURL(secret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", rec.Code)
		}
	})

	t.Run("expired link returns 403", func(t *testing.T) {
		// Generate a URL that expired 10 minutes ago.
		url := GenerateSignedURL("http://localhost:8080/api/v1/user/email/verify", 42, secret, -10)
		req := httptest.NewRequest(http.MethodGet, url, nil)
		rec := httptest.NewRecorder()

		handler := SignedURL(secret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d", rec.Code)
		}
	})

	t.Run("wrong signature returns 401", func(t *testing.T) {
		url := GenerateSignedURL("http://localhost:8080/api/v1/user/email/verify", 42, secret, 60)
		// Tamper with the signature.
		url = strings.Replace(url, "signature=", "signature=tampered", 1)
		req := httptest.NewRequest(http.MethodGet, url, nil)
		rec := httptest.NewRecorder()

		handler := SignedURL(secret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", rec.Code)
		}
	})
}

func TestComputeHMAC(t *testing.T) {
	t.Run("deterministic", func(t *testing.T) {
		h1 := computeHMAC("test-data", "secret")
		h2 := computeHMAC("test-data", "secret")
		if h1 != h2 {
			t.Fatal("same input should produce same HMAC")
		}
	})

	t.Run("different data produces different HMAC", func(t *testing.T) {
		h1 := computeHMAC("data1", "secret")
		h2 := computeHMAC("data2", "secret")
		if h1 == h2 {
			t.Fatal("different data should produce different HMAC")
		}
	})

	t.Run("different secret produces different HMAC", func(t *testing.T) {
		h1 := computeHMAC("data", "secret1")
		h2 := computeHMAC("data", "secret2")
		if h1 == h2 {
			t.Fatal("different secrets should produce different HMAC")
		}
	})
}
