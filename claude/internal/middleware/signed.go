package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/incharge/server/internal/dto"
)

// GenerateSignedURL generates a signed URL for email verification.
func GenerateSignedURL(baseURL string, userID uint, secret string, expiryMinutes int) string {
	expires := time.Now().Add(time.Duration(expiryMinutes) * time.Minute).Unix()
	// Sign the full path + query (not the full URL), matching the verifier.
	data := fmt.Sprintf("%s/%d?expires=%d", baseURL, userID, expires)
	signature := computeHMAC(data, secret)
	return fmt.Sprintf("%s/%d?expires=%d&signature=%s", baseURL, userID, expires, signature)
}

// SignedURL validates that the request URL has a valid signature and has not expired.
func SignedURL(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			signature := r.URL.Query().Get("signature")
			expiresStr := r.URL.Query().Get("expires")

			if signature == "" || expiresStr == "" {
				dto.WriteAuthError(w, "Invalid or missing signature")
				return
			}

			expires, err := strconv.ParseInt(expiresStr, 10, 64)
			if err != nil {
				dto.WriteAuthError(w, "Invalid expiry")
				return
			}

			if time.Now().Unix() > expires {
				dto.WriteJSON(w, http.StatusForbidden, map[string]string{
					"message": "This link has expired.",
				})
				return
			}

			// Reconstruct the full URL that was signed (scheme + host + path + query).
			// GenerateSignedURL signs: baseURL/id?expires=...
			// We reconstruct using the request's scheme/host/path.
			scheme := "http"
			if r.TLS != nil {
				scheme = "https"
			}
			if fwd := r.Header.Get("X-Forwarded-Proto"); fwd != "" {
				scheme = fwd
			}
			fullURL := fmt.Sprintf("%s://%s%s", scheme, r.Host, r.URL.Path)
			data := fmt.Sprintf("%s?expires=%d", fullURL, expires)
			expected := computeHMAC(data, secret)

			if !hmac.Equal([]byte(signature), []byte(expected)) {
				dto.WriteAuthError(w, "Invalid signature")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func computeHMAC(data, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
