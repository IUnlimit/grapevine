package config

import (
	"os"
	"strings"
)

type Config struct {
	Port      string
	DBURL     string
	LogLevel  string
	JWTSecret string
}

func Load() *Config {
	return &Config{
		Port:      getEnv("GRAPEVINE_PORT", "8888"),
		DBURL:     getEnv("DB_URL", "postgres://postgres:123456@127.0.0.1:5432/grapevine?sslmode=disable"),
		LogLevel:  strings.ToLower(getEnv("LOG_LEVEL", "info")),
		JWTSecret: getEnv("JWT_SECRET", "grapevine_secret"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
