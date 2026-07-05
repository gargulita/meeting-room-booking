package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	DBMaxOpenConns  int
	DBMaxIdleConns  int
	Port            string
	JWTSecret       string
	SlotsLookupDays int
}

func Load() *Config {
	return &Config{
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "postgres"),
		DBPassword:      getEnv("DB_PASSWORD", "password"),
		DBName:          getEnv("DB_NAME", "booking_service"),
		DBMaxOpenConns:  getEnvInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:  getEnvInt("DB_MAX_IDLE_CONNS", 10),
		Port:            getEnv("PORT", "8080"),
		JWTSecret:       getEnv("JWT_SECRET", "your-secret-key"),
		SlotsLookupDays: getEnvInt("SLOTS_LOOKUP_DAYS", 7),
	}
}

func (c *Config) DBConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intVal int
		fmt.Sscanf(value, "%d", &intVal)
		return intVal
	}
	return defaultValue
}
