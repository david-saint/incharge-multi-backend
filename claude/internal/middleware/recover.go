package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/incharge/server/internal/dto"
)

// Recover catches panics and returns a structured 500 response.
func Recover(isProduction bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					slog.Error("panic recovered",
						"error", rec,
						"path", r.URL.Path,
						"method", r.Method,
					)

					resp := dto.ErrorResponse{
						Status:  false,
						Message: "Internal server error",
					}
					if !isProduction {
						resp.Error = fmt.Sprintf("%v", rec)
						resp.Trace = string(debug.Stack())
					}
					dto.WriteJSON(w, http.StatusInternalServerError, resp)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
