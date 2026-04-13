package database

import (
	"fmt"
	"log"

	"incharge/internal/config"
	"incharge/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg config.Config) {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUsername,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBDatabase,
	)

	gormConfig := &gorm.Config{}
	if cfg.AppEnv == "local" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	} else {
		gormConfig.Logger = logger.Default.LogMode(logger.Error)
	}

	DB, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established")

	// Auto-migrate models
	err = DB.AutoMigrate(
		&models.User{},
		&models.Profile{},
		&models.Clinic{},
		&models.Location{},
		&models.State{},
		&models.Country{},
		&models.ContraceptionReason{},
		&models.EducationLevel{},
		&models.FaqGroup{},
		&models.Faq{},
		&models.Algorithm{},
		&models.Admin{},
		&models.PasswordReset{},
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}
}
