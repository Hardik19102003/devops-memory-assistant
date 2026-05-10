package config

import "os"

type Config struct {
	DBURL string
	PORT  string
	ENV   string
}

func Load() Config {
	return Config{
		DBURL: getEnv("DB_URL", "postgres://localhost:5432/dev"),
		PORT:  getEnv("PORT", "8080"),
		ENV:   getEnv("ENV", "dev"),
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}