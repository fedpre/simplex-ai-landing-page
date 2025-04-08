package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	AdminEmail     string
	AdminPassword  string
	Environment    string
	CookieDuration time.Duration
}

func LoadConfig() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// Debug logging for environment variables
	adminEmail := getEnv("ADMIN_EMAIL", "admin")
	adminPassword := getEnv("ADMIN_PASSWORD", "admin")

	return &Config{
		Port:           getEnv("PORT", "8080"),
		AdminEmail:     adminEmail,
		AdminPassword:  adminPassword,
		Environment:    getEnv("ENVIRONMENT", "development"),
		CookieDuration: time.Hour,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
