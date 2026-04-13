package database

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/incharge/server/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect establishes a MySQL connection and returns a configured GORM instance.
func Connect(cfg *config.DBConfig) (*gorm.DB, error) {
	logLevel := logger.Warn
	if !isProduction(cfg.Host) {
		logLevel = logger.Info
	}

	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger:                 logger.Default.LogMode(logLevel),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(2 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	slog.Info("database connected", "host", cfg.Host, "database", cfg.Database)
	return db, nil
}

func isProduction(host string) bool {
	return host != "127.0.0.1" && host != "localhost"
}
