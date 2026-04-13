package main

import (
	"log"

	"incharge/internal/api/routes"
	"incharge/internal/config"
	"incharge/internal/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	// Initialize Database
	database.InitDB(cfg)

	// Set Gin mode
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	r := gin.Default()

	// Session Setup
	store := cookie.NewStore([]byte(cfg.JWTSecret)) // Reusing JWT secret for session encryption
	r.Use(sessions.Sessions("incharge_session", store))

	// Apply CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Authorization"},
		AllowCredentials: true,
	}))

	routes.SetupRoutes(r, cfg)

	// Run Server
	log.Printf("Starting server on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
