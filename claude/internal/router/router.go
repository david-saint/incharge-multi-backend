package router

import (
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/incharge/server/internal/config"
	"github.com/incharge/server/internal/handler"
	"github.com/incharge/server/internal/middleware"
	"github.com/incharge/server/internal/repository"
	"github.com/incharge/server/internal/service"
	"gorm.io/gorm"
)

// Setup creates and configures the Chi router with all routes.
func Setup(cfg *config.Config, db *gorm.DB) *chi.Mux {
	r := chi.NewRouter()

	// --- Global middleware ---
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(middleware.Recover(cfg.App.IsProduction))
	r.Use(middleware.CORS())

	// Rate limiters.
	globalLimiter := middleware.NewRateLimiter(120, time.Minute)
	emailLimiter := middleware.NewRateLimiter(6, time.Minute)
	r.Use(globalLimiter.Handler)

	// Initialize session store.
	middleware.InitSessionStore(cfg.Session.Secret)

	// --- Create repositories ---
	userRepo := repository.NewUserRepo(db)
	profileRepo := repository.NewProfileRepo(db)
	clinicRepo := repository.NewClinicRepo(db)
	adminRepo := repository.NewAdminRepo(db)
	algoRepo := repository.NewAlgorithmRepo(db)
	refRepo := repository.NewReferenceRepo(db)

	// --- Create services ---
	authService := service.NewAuthService(cfg.JWT.Secret, cfg.JWT.TTL)
	emailService := service.NewEmailService(&cfg.Mail)

	// --- Create handlers ---
	authHandler := handler.NewAuthHandler(userRepo, authService, emailService, cfg)
	globalHandler := handler.NewGlobalHandler(refRepo, cfg)
	profileHandler := handler.NewProfileHandler(profileRepo, cfg)
	clinicHandler := handler.NewClinicHandler(clinicRepo, cfg)
	userHandler := handler.NewUserHandler(userRepo, cfg)
	adminHandler := handler.NewAdminHandler(adminRepo, userRepo, clinicRepo, algoRepo, refRepo, authService, cfg)
	algoHandler := handler.NewAlgorithmHandler(algoRepo, cfg)

	// =========================================================================
	// API v1 Routes
	// =========================================================================
	r.Route("/api/v1", func(r chi.Router) {
		// --- Global (public) routes ---
		r.Route("/global", func(r chi.Router) {
			r.Get("/", globalHandler.HealthCheck) // Health check at /api/v1/global/
			r.Get("/contraception-reasons", globalHandler.ListContraceptionReasons)
			r.Get("/contraception-reasons/{id}", globalHandler.GetContraceptionReason)
			r.Get("/education-levels", globalHandler.ListEducationLevels)
			r.Get("/faq-groups", globalHandler.ListFaqGroups)
			r.Get("/faq-groups/{id}", globalHandler.GetFaqGroup)
		})

		// --- User auth routes ---
		r.Route("/user", func(r chi.Router) {
			// Public auth routes.
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
			r.Post("/logout", authHandler.Logout)
			r.Get("/refresh", authHandler.Refresh)

			// Password reset.
			r.Post("/password/email", authHandler.PasswordEmail)
			r.Post("/password/reset", authHandler.PasswordReset)

			// Email verification.
			r.With(middleware.SignedURL(cfg.JWT.Secret)).
				Get("/email/verify/{id}", authHandler.VerifyEmail)
			r.With(middleware.JWTAuth(cfg.JWT.Secret), emailLimiter.Handler).
				Get("/email/resend", authHandler.ResendVerification)
			r.Get("/email/success", authHandler.EmailSuccess)

			// Authenticated user endpoint.
			r.With(middleware.JWTAuth(cfg.JWT.Secret)).Get("/", authHandler.GetUser)

			// --- Profile routes (JWT protected) ---
			r.Route("/profile", func(r chi.Router) {
				r.Use(middleware.JWTAuth(cfg.JWT.Secret))
				r.Post("/", profileHandler.SaveProfile)
				r.Get("/", profileHandler.GetProfile)
				r.Post("/algorithm", profileHandler.StoreAlgorithmPlan)
			})

			// --- Clinic routes ---
			r.Route("/clinics", func(r chi.Router) {
				r.Get("/", clinicHandler.ListClinics) // Public (no auth).

				// Session-protected clinic management.
				r.With(middleware.IsLogged).Get("/getClinics", clinicHandler.GetClinicsPaginated)
				r.With(middleware.IsLogged).Get("/deletedClinics", clinicHandler.GetDeletedClinics)
				r.With(middleware.IsLogged).Post("/addClinic", clinicHandler.AddClinic)
				r.With(middleware.IsLogged).Put("/update/{clinic_id}", clinicHandler.UpdateClinic)
				r.With(middleware.IsLogged).Put("/revertDelete/{clinic_id}", clinicHandler.RevertDeleteClinic)
				r.With(middleware.IsLogged).Delete("/deleteClinic/{clinic_id}", clinicHandler.DeleteClinic)
			})

			// --- User management routes (session protected) ---
			r.Route("/users", func(r chi.Router) {
				r.Use(middleware.IsLogged)
				r.Get("/", userHandler.ListUsers)
				r.Get("/deletedUser", userHandler.ListDeletedUsers)
				r.Put("/update/{user_id}", userHandler.RestoreUser)
				r.Delete("/deleteUser/{user_id}", userHandler.DeleteUser)
			})
		})
	})

	// =========================================================================
	// Admin Web Panel Routes
	// =========================================================================

	// Public admin web routes.
	r.Get("/", adminHandler.Index)
	r.Get("/loginView", adminHandler.LoginView)
	r.Post("/login", adminHandler.Login)
	r.Get("/logout", adminHandler.Logout)
	r.Get("/privacy", adminHandler.Privacy)

	// Algorithm list is publicly accessible.
	r.Get("/algo", algoHandler.ListAlgorithms)

	// POST /admin — accessible without session when no Super admin exists (first-time setup).
	// The handler itself checks whether auth is required.
	r.Post("/admin", adminHandler.CreateAdmin)

	// Protected admin routes (require verified admin session).
	r.Group(func(r chi.Router) {
		r.Use(middleware.IsAdmin(db))

		// Admin panel.
		r.Get("/panel", adminHandler.Panel)

		// Admin management.
		r.Get("/allAdmins", adminHandler.ListAdmins)
		r.Get("/getAdminDet", adminHandler.GetAdminDetails)
		r.Put("/updateAdmin/{admin_id}", adminHandler.UpdateAdmin)

		// Algorithm management.
		r.Post("/algo", algoHandler.CreateAlgorithm)
		r.Put("/algo/{id}", algoHandler.UpdateAlgorithm)

		// User management.
		r.Get("/getUsers", adminHandler.AdminListUsers)
		r.Put("/updateUser/{user_id}", adminHandler.AdminRestoreUser)
		r.Delete("/deleteUser/{user_id}", adminHandler.AdminDeleteUser)
		r.Get("/getDeletedUsers", adminHandler.AdminListDeletedUsers)
		r.Put("/revertDeletedUser/{user_id}", adminHandler.AdminRestoreUser)

		// Clinic management.
		r.Get("/getClinics", adminHandler.AdminListClinics)
		r.Post("/addClinic", adminHandler.AdminAddClinic)
		r.Put("/updateClinic/{clinic_id}", adminHandler.AdminUpdateClinic)
		r.Delete("/deleteClinic/{clinic_id}", adminHandler.AdminDeleteClinic)
		r.Get("/getDeletedClinics", adminHandler.AdminListDeletedClinics)
		r.Put("/revertDeletedClinic/{clinic_id}", adminHandler.AdminRestoreClinic)

		// Reference data.
		r.Get("/getContraceptionReason", adminHandler.AdminListContraceptionReasons)
		r.Get("/getEducationalLevels", adminHandler.AdminListEducationLevels)
	})

	return r
}
