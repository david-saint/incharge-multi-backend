package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

// CORS returns a CORS middleware configured per spec:
// allow all origins, all headers, all methods; expose Authorization header.
func CORS() func(http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Authorization"},
		AllowCredentials: false,
		MaxAge:           86400,
	})
	return c.Handler
}
