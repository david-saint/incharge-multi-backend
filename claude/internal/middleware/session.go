package middleware

import (
	"context"
	"encoding/gob"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/incharge/server/internal/dto"
	"gorm.io/gorm"

	"github.com/incharge/server/internal/model"
)

const (
	sessionName     = "incharge_session"
	sessionAdminKey = "admin_id"
)

// SessionStore is the global session store.
var SessionStore *sessions.CookieStore

// InitSessionStore creates the session cookie store.
func InitSessionStore(secret string) {
	// Register uint type with gob so gorilla/sessions can encode/decode it.
	gob.Register(uint(0))

	SessionStore = sessions.NewCookieStore([]byte(secret))
	SessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

// IsLogged checks that the session has an authenticated admin.
func IsLogged(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := SessionStore.Get(r, sessionName)
		if err != nil || session.IsNew {
			dto.WriteAuthError(w, "Permission Denied")
			return
		}

		adminID, ok := session.Values[sessionAdminKey].(uint)
		if !ok || adminID == 0 {
			dto.WriteAuthError(w, "Permission Denied")
			return
		}

		ctx := context.WithValue(r.Context(), AdminIDKey, adminID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// IsAdmin checks that the session has a verified admin.
func IsAdmin(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := SessionStore.Get(r, sessionName)
			if err != nil || session.IsNew {
				dto.WriteAuthError(w, "Permission Denied")
				return
			}

			adminID, ok := session.Values[sessionAdminKey].(uint)
			if !ok || adminID == 0 {
				dto.WriteAuthError(w, "Permission Denied")
				return
			}

			var admin model.Admin
			if err := db.First(&admin, adminID).Error; err != nil {
				dto.WriteAuthError(w, "Permission Denied")
				return
			}

			if !admin.IsVerified() {
				dto.WriteAuthError(w, "Permission Denied")
				return
			}

			ctx := context.WithValue(r.Context(), AdminIDKey, adminID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// SetAdminSession stores the admin ID in the session.
func SetAdminSession(w http.ResponseWriter, r *http.Request, adminID uint) error {
	session, err := SessionStore.Get(r, sessionName)
	if err != nil {
		session, _ = SessionStore.New(r, sessionName)
	}
	session.Values[sessionAdminKey] = adminID
	return session.Save(r, w)
}

// ClearAdminSession invalidates the admin session.
func ClearAdminSession(w http.ResponseWriter, r *http.Request) error {
	session, err := SessionStore.Get(r, sessionName)
	if err != nil {
		return err
	}
	session.Values[sessionAdminKey] = uint(0)
	session.Options.MaxAge = -1
	return session.Save(r, w)
}
