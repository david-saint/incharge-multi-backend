package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/incharge/server/internal/dto"
)

// rateLimitEntry tracks request counts for an IP.
type rateLimitEntry struct {
	count       int
	windowStart time.Time
}

// RateLimiter implements a sliding-window rate limiter.
type RateLimiter struct {
	mu      sync.Mutex
	entries map[string]*rateLimitEntry
	limit   int
	window  time.Duration
}

// NewRateLimiter creates a new rate limiter.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		entries: make(map[string]*rateLimitEntry),
		limit:   limit,
		window:  window,
	}
	// Cleanup goroutine.
	go func() {
		ticker := time.NewTicker(window)
		defer ticker.Stop()
		for range ticker.C {
			rl.cleanup()
		}
	}()
	return rl
}

// Handler returns an HTTP middleware that enforces the rate limit.
func (rl *RateLimiter) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := realIP(r)

		rl.mu.Lock()
		entry, exists := rl.entries[ip]
		now := time.Now()
		if !exists || now.Sub(entry.windowStart) > rl.window {
			rl.entries[ip] = &rateLimitEntry{count: 1, windowStart: now}
			rl.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}
		entry.count++
		if entry.count > rl.limit {
			rl.mu.Unlock()
			dto.WriteJSON(w, http.StatusTooManyRequests, map[string]string{
				"message": "Too many requests. Please try again later.",
			})
			return
		}
		rl.mu.Unlock()
		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	for ip, entry := range rl.entries {
		if now.Sub(entry.windowStart) > rl.window {
			delete(rl.entries, ip)
		}
	}
}

func realIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	return r.RemoteAddr
}
