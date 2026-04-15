package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	AppEnv         string `mapstructure:"APP_ENV"`
	AppUrl         string `mapstructure:"APP_URL"`
	ApiDomain      string `mapstructure:"app.api-domain"`
	UserDomain     string `mapstructure:"app.user-domain"`
	DBConnection   string `mapstructure:"DB_CONNECTION"`
	DBHost         string `mapstructure:"DB_HOST"`
	DBPort         string `mapstructure:"DB_PORT"`
	DBDatabase     string `mapstructure:"DB_DATABASE"`
	DBUsername     string `mapstructure:"DB_USERNAME"`
	DBPassword     string `mapstructure:"DB_PASSWORD"`
	JWTSecret      string `mapstructure:"JWT_SECRET"`
	MailHost       string `mapstructure:"MAIL_HOST"`
	MailPort       string `mapstructure:"MAIL_PORT"`
	MailUsername   string `mapstructure:"MAIL_USERNAME"`
	MailPassword   string `mapstructure:"MAIL_PASSWORD"`
	MailEncryption string `mapstructure:"MAIL_ENCRYPTION"`
	Port           string `mapstructure:"PORT"`
}

func LoadConfig(path string) (config Config, err error) {
	godotenv.Load(path + "/.env") // Load .env file if present

	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
		err = nil // we don't care if the file isn't there, as long as ENV vars exist
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	config.Port = getString("PORT", config.Port, "8080")
	config.AppEnv = getString("APP_ENV", config.AppEnv, "local")
	config.AppUrl = getString("APP_URL", config.AppUrl, "http://localhost:"+config.Port)
	config.ApiDomain = getString("APP_API_DOMAIN", config.ApiDomain, config.AppUrl)
	config.UserDomain = getString("APP_USER_DOMAIN", config.UserDomain, "http://localhost:3000")
	config.DBConnection = getString("DB_CONNECTION", config.DBConnection, "mysql")
	config.DBHost = getString("DB_HOST", config.DBHost, "127.0.0.1")
	config.DBPort = getString("DB_PORT", config.DBPort, "3306")
	config.DBDatabase = getString("DB_DATABASE", config.DBDatabase, "incharge")
	config.DBUsername = getString("DB_USERNAME", config.DBUsername, "incharge")
	config.DBPassword = getString("DB_PASSWORD", config.DBPassword, "secret")
	config.JWTSecret = getString("JWT_SECRET", config.JWTSecret, "change-me-secret")
	config.MailHost = getString("MAIL_HOST", config.MailHost, "")
	config.MailPort = getString("MAIL_PORT", config.MailPort, "1025")
	config.MailUsername = getString("MAIL_USERNAME", config.MailUsername, "")
	config.MailPassword = getString("MAIL_PASSWORD", config.MailPassword, "")
	config.MailEncryption = getString("MAIL_ENCRYPTION", config.MailEncryption, "")

	return
}

func getString(key, currentValue, fallback string) string {
	if value := viper.GetString(key); value != "" {
		return value
	}
	if currentValue != "" {
		return currentValue
	}
	return fallback
}
