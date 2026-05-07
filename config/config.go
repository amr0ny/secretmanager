package config

import "os"

type Config struct {
	DatabaseDSN   string
	HashSecretKey string
}

func Load() (*Config, error) {
	return &Config{
		DatabaseDSN:   getEnv("DATABASE_DSN", ""),
		HashSecretKey: getEnv("HASH_SECRET_KEY", ""),
	}, nil
}
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
