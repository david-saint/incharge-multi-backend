package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/incharge/server/internal/config"
	"github.com/incharge/server/internal/database"
	"github.com/incharge/server/internal/router"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists (non-fatal if missing).
	if err := godotenv.Load(); err != nil {
		slog.Info("no .env file found, using environment variables")
	}

	// Configure structured logging.
	logLevel := slog.LevelInfo
	if os.Getenv("APP_ENV") == "local" {
		logLevel = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})))

	// Load configuration.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	slog.Info("starting InCharge server",
		"env", cfg.App.Env,
		"port", cfg.App.Port,
	)

	// Connect to database.
	db, err := database.Connect(&cfg.DB)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	// Run auto-migration.
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("auto-migration failed: %v", err)
	}

	// Seed reference data.
	if err := database.Seed(db); err != nil {
		log.Fatalf("seeding failed: %v", err)
	}

	// Setup router.
	r := router.Setup(cfg, db)

	// Configure HTTP server.
	srv := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	slog.Info("server listening", "addr", srv.Addr, "url", cfg.App.URL)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
