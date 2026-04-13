package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	t.Run("allows requests under limit", func(t *testing.T) {
		rl := NewRateLimiter(5, 1*time.Minute)

		handler := rl.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		for i := 0; i < 5; i++ {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = "192.168.1.1:12345"
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			if rec.Code != http.StatusOK {
				t.Fatalf("request %d: expected 200, got %d", i+1, rec.Code)
			}
		}
	})

	t.Run("blocks requests over limit", func(t *testing.T) {
		rl := NewRateLimiter(3, 1*time.Minute)

		handler := rl.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		for i := 0; i < 3; i++ {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = "10.0.0.1:12345"
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
		}

		// 4th request should be blocked.
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "10.0.0.1:12345"
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusTooManyRequests {
			t.Fatalf("expected 429, got %d", rec.Code)
		}
	})

	t.Run("different IPs are independent", func(t *testing.T) {
		rl := NewRateLimiter(2, 1*time.Minute)

		handler := rl.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		// Exhaust limit for IP A.
		for i := 0; i < 2; i++ {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = "1.1.1.1:123"
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
		}

		// IP B should still be allowed.
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "2.2.2.2:123"
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 for different IP, got %d", rec.Code)
		}
	})

	t.Run("uses X-Forwarded-For header", func(t *testing.T) {
		rl := NewRateLimiter(1, 1*time.Minute)

		handler := rl.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		// First request from forwarded IP.
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "127.0.0.1:123"
		req.Header.Set("X-Forwarded-For", "3.3.3.3")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("first request should pass, got %d", rec.Code)
		}

		// Second request from same forwarded IP should be blocked.
		req2 := httptest.NewRequest(http.MethodGet, "/", nil)
		req2.RemoteAddr = "127.0.0.1:456"
		req2.Header.Set("X-Forwarded-For", "3.3.3.3")
		rec2 := httptest.NewRecorder()
		handler.ServeHTTP(rec2, req2)
		if rec2.Code != http.StatusTooManyRequests {
			t.Fatalf("second request should be blocked, got %d", rec2.Code)
		}
	})
}

func TestRealIP(t *testing.T) {
	t.Run("prefers X-Forwarded-For", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "127.0.0.1:123"
		req.Header.Set("X-Forwarded-For", "1.2.3.4")
		req.Header.Set("X-Real-IP", "5.6.7.8")
		if ip := realIP(req); ip != "1.2.3.4" {
			t.Fatalf("expected 1.2.3.4, got %s", ip)
		}
	})

	t.Run("falls back to X-Real-IP", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "127.0.0.1:123"
		req.Header.Set("X-Real-IP", "5.6.7.8")
		if ip := realIP(req); ip != "5.6.7.8" {
			t.Fatalf("expected 5.6.7.8, got %s", ip)
		}
	})

	t.Run("falls back to RemoteAddr", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "127.0.0.1:123"
		if ip := realIP(req); ip != "127.0.0.1:123" {
			t.Fatalf("expected 127.0.0.1:123, got %s", ip)
		}
	})
}
