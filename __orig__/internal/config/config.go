package config

import (
	"os"
	"time"
)

type Config struct {
	Addr           string
	DBPath         string
	WorkerInterval time.Duration
	SessionTTL     time.Duration
}

func Load() Config {
	return Config{Addr: env("ADDR", ":8080"), DBPath: env("DB_PATH", "./data/lanzhou.db"), WorkerInterval: duration("WORKER_INTERVAL", 30*time.Second), SessionTTL: duration("SESSION_TTL", 8*time.Hour)}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
func duration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}
	return fallback
}
