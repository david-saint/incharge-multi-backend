package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration.
type Config struct {
	App     AppConfig
	DB      DBConfig
	JWT     JWTConfig
	Mail    MailConfig
	Session SessionConfig
}

// AppConfig holds application-level settings.
type AppConfig struct {
	Env          string
	Port         string
	URL          string
	APIDomain    string
	UserDomain   string
	IsProduction bool
}

// DBConfig holds database connection settings.
type DBConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
	DSN      string
}

// JWTConfig holds JWT signing settings.
type JWTConfig struct {
	Secret string
	TTL    time.Duration
}

// MailConfig holds SMTP settings.
type MailConfig struct {
	Host       string
	Port       int
	Username   string
	Password   string
	Encryption string
	FromAddr   string
	FromName   string
}

// SessionConfig holds session cookie settings.
type SessionConfig struct {
	Secret string
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	appEnv := getEnv("APP_ENV", "local")
	port := getEnv("APP_PORT", "8080")
	jwtTTL := getEnvAsInt("JWT_TTL_MINUTES", 1440)
	mailPort := getEnvAsInt("MAIL_PORT", 587)

	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_DATABASE", "")
	dbUser := getEnv("DB_USERNAME", "")
	dbPass := getEnv("DB_PASSWORD", "")

	if dbName == "" {
		return nil, fmt.Errorf("DB_DATABASE is required")
	}
	if dbUser == "" {
		return nil, fmt.Errorf("DB_USERNAME is required")
	}

	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	return &Config{
		App: AppConfig{
			Env:          appEnv,
			Port:         port,
			URL:          getEnv("APP_URL", "http://localhost:"+port),
			APIDomain:    getEnv("APP_API_DOMAIN", "localhost"),
			UserDomain:   getEnv("APP_USER_DOMAIN", "http://localhost:3000"),
			IsProduction: appEnv == "production",
		},
		DB: DBConfig{
			Host:     dbHost,
			Port:     dbPort,
			Database: dbName,
			Username: dbUser,
			Password: dbPass,
			DSN:      dsn,
		},
		JWT: JWTConfig{
			Secret: jwtSecret,
			TTL:    time.Duration(jwtTTL) * time.Minute,
		},
		Mail: MailConfig{
			Host:       getEnv("MAIL_HOST", ""),
			Port:       mailPort,
			Username:   getEnv("MAIL_USERNAME", ""),
			Password:   getEnv("MAIL_PASSWORD", ""),
			Encryption: getEnv("MAIL_ENCRYPTION", "tls"),
			FromAddr:   getEnv("MAIL_FROM_ADDRESS", "noreply@incharge.app"),
			FromName:   getEnv("MAIL_FROM_NAME", "InCharge"),
		},
		Session: SessionConfig{
			Secret: getEnv("APP_SESSION_SECRET", "incharge-session-secret-change-me"),
		},
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return i
}
