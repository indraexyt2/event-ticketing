package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppName     string
	AppPort     string
	LogLevel    string
	Environment string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTSecret    string
	JWTExpiresIn time.Duration
}

func LoadConfig() Config {
	config := Config{
		AppName:     getEnv("APP_NAME", "Ticketing System API"),
		AppPort:     getEnv("APP_PORT", "8080"),
		LogLevel:    getEnv("LOG_LEVEL", "debug"),
		Environment: getEnv("ENVIRONMENT", "development"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "ticketing_system"),

		JWTSecret:    getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpiresIn: time.Duration(getEnvAsInt("JWT_EXPIRES_IN", 24)) * time.Hour,
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
