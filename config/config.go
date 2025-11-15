package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                 string
	Domain               string
	DatabaseUrl          string
	JWTSecret            string
	JWTAccessExpiration  string
	JWTRefreshExpiration string
	EmailsAPIKey         string
	SenderEmail          string
}

func LoadConfig() (*Config, error) {
	godotenv.Load()

	domain := getEnv("PROD_URL", "")
	if getEnv("DEV_MODE", "") == "dev" {
		domain = getEnv("DEV_URL", "")
	}

	return &Config{
		Port:                 getEnv("PORT", "8080"),
		Domain:               domain,
		DatabaseUrl:          getEnv("DATABASE_URL", ""),
		JWTSecret:            getEnv("JWT_SECRET", "super-secret-key"),
		JWTAccessExpiration:  getEnv("JWT_ACCESS_EXPIRATION", "1h"),
		JWTRefreshExpiration: getEnv("JWT_REFRESH_EXPIRATION", "7d"),
		EmailsAPIKey:         getEnv("EMAILS_API_KEY", ""),
		SenderEmail:          getEnv("SENDER_EMAIL", ""),
	}, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
