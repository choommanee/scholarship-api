package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL      string
	JWTSecret        string
	Port             string
	Environment      string
	UploadPath       string
	MaxFileSize      int64
	SSOEnabled       bool
	SSOEntityID      string
	SSOSSOUrl        string
	SSOCertPath      string
	EmailSMTPHost    string
	EmailSMTPPort    string
	EmailUsername    string
	EmailPassword    string
	EmailFrom        string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	return &Config{
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/scholarship_db?sslmode=disable"),
		JWTSecret:        getEnv("JWT_SECRET", "your-secret-key"),
		Port:             getEnv("PORT", "8080"),
		Environment:      getEnv("ENVIRONMENT", "development"),
		UploadPath:       getEnv("UPLOAD_PATH", "./uploads"),
		MaxFileSize:      getEnvInt64("MAX_FILE_SIZE", 10485760), // 10MB
		SSOEnabled:       getEnv("SSO_ENABLED", "false") == "true",
		SSOEntityID:      getEnv("SSO_ENTITY_ID", ""),
		SSOSSOUrl:        getEnv("SSO_SSO_URL", ""),
		SSOCertPath:      getEnv("SSO_CERT_PATH", ""),
		EmailSMTPHost:    getEnv("EMAIL_SMTP_HOST", "localhost"),
		EmailSMTPPort:    getEnv("EMAIL_SMTP_PORT", "587"),
		EmailUsername:    getEnv("EMAIL_USERNAME", ""),
		EmailPassword:    getEnv("EMAIL_PASSWORD", ""),
		EmailFrom:        getEnv("EMAIL_FROM", "noreply@university.ac.th"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}