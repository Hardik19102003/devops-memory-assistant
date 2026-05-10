package config

import "os"

type Config struct {
	DBURL string
	PORT  string
	ENV   string
}

func Load() Config {
	return Config{
		DBURL: getEnv("DB_URL", "postgres://devops:devops@localhost:5432/devops_memory?sslmode=disable"),
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