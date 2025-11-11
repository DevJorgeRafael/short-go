package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                 string
	DBPath               string
	JWTSecret            string
	JWTAccessExpiration  string
	JWTRefreshExpiration string
}

func LoadConfig() (*Config, error) {
	godotenv.Load()

	return &Config{
		Port: getEnv("PORT", "8080"),
		DBPath: getEnv("DB_PATH", "./todo.db"),
		JWTSecret: getEnv("JWT_SECRET", "super-secret-key"),
		JWTAccessExpiration: getEnv("JWT_ACCESS_EXPIRATION", "1h"),
		JWTRefreshExpiration: getEnv("JWT_REFRESH_EXPIRATION", "7d"),
	} , nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
