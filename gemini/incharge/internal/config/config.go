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

	// Defaults
	if config.Port == "" {
		config.Port = "8080"
	}
	if config.AppEnv == "" {
		config.AppEnv = "local"
	}

	return
}
