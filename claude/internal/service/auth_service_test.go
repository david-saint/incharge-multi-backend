package service

import (
	"testing"
	"time"
)

func TestGenerateAndValidateToken(t *testing.T) {
	svc := NewAuthService("test-secret-key-12345", 60*time.Minute)

	t.Run("valid token round-trip", func(t *testing.T) {
		token, err := svc.GenerateToken(42)
		if err != nil {
			t.Fatalf("GenerateToken: %v", err)
		}
		if token == "" {
			t.Fatal("token should not be empty")
		}

		userID, err := svc.ValidateToken(token)
		if err != nil {
			t.Fatalf("ValidateToken: %v", err)
		}
		if userID != 42 {
			t.Fatalf("expected userID=42, got %d", userID)
		}
	})

	t.Run("different user IDs produce different tokens", func(t *testing.T) {
		t1, _ := svc.GenerateToken(1)
		t2, _ := svc.GenerateToken(2)
		if t1 == t2 {
			t.Fatal("tokens for different users should differ")
		}
	})

	t.Run("invalid token string", func(t *testing.T) {
		_, err := svc.ValidateToken("garbage.token.here")
		if err == nil {
			t.Fatal("expected error for invalid token")
		}
	})

	t.Run("wrong secret rejects token", func(t *testing.T) {
		token, _ := svc.GenerateToken(1)
		other := NewAuthService("wrong-secret", 60*time.Minute)
		_, err := other.ValidateToken(token)
		if err == nil {
			t.Fatal("expected error with wrong secret")
		}
	})

	t.Run("expired token is rejected", func(t *testing.T) {
		expired := NewAuthService("test-secret-key-12345", -1*time.Minute)
		token, _ := expired.GenerateToken(1)
		_, err := svc.ValidateToken(token)
		if err == nil {
			t.Fatal("expected error for expired token")
		}
	})
}

func TestHashAndCheckPassword(t *testing.T) {
	svc := NewAuthService("unused", 60*time.Minute)

	t.Run("correct password", func(t *testing.T) {
		hash, err := svc.HashPassword("mypassword123")
		if err != nil {
			t.Fatalf("HashPassword: %v", err)
		}
		if !svc.CheckPassword(hash, "mypassword123") {
			t.Fatal("expected correct password to match")
		}
	})

	t.Run("wrong password", func(t *testing.T) {
		hash, _ := svc.HashPassword("mypassword123")
		if svc.CheckPassword(hash, "wrongpassword") {
			t.Fatal("expected wrong password to not match")
		}
	})

	t.Run("different hashes for same password", func(t *testing.T) {
		h1, _ := svc.HashPassword("same")
		h2, _ := svc.HashPassword("same")
		if h1 == h2 {
			t.Fatal("bcrypt should produce different hashes (salted)")
		}
	})
}

func TestGenerateRandomToken(t *testing.T) {
	t.Run("produces hex token of expected length", func(t *testing.T) {
		token, err := GenerateRandomToken(32)
		if err != nil {
			t.Fatalf("GenerateRandomToken: %v", err)
		}
		if len(token) != 64 { // 32 bytes = 64 hex chars
			t.Fatalf("expected 64 hex chars, got %d", len(token))
		}
	})

	t.Run("produces unique tokens", func(t *testing.T) {
		t1, _ := GenerateRandomToken(16)
		t2, _ := GenerateRandomToken(16)
		if t1 == t2 {
			t.Fatal("random tokens should be unique")
		}
	})
}
