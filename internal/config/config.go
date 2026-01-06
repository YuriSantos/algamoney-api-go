package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	GinMode         string
	DBHost          string
	DBPort          string
	DBDatabase      string
	DBUsername      string
	DBPassword      string
	JWTSecret       string
	JWTExpiresHours int
	CORSOrigin      string
}

func Load() *Config {
	godotenv.Load()

	jwtExpires, _ := strconv.Atoi(getEnv("JWT_EXPIRES_HOURS", "24"))

	return &Config{
		Port:            getEnv("PORT", "8080"),
		GinMode:         getEnv("GIN_MODE", "debug"),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "3306"),
		DBDatabase:      getEnv("DB_DATABASE", "algamoneyapi"),
		DBUsername:      getEnv("DB_USERNAME", "root"),
		DBPassword:      getEnv("DB_PASSWORD", ""),
		JWTSecret:       getEnv("JWT_SECRET", "secret"),
		JWTExpiresHours: jwtExpires,
		CORSOrigin:      getEnv("CORS_ORIGIN", "http://localhost:4200"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
