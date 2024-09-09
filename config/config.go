package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser                 string
	DBPassword             string
	DBHost                 string
	DBName                 string
	PETDB                  string
	JWTSecret              string
	JWTExpirationInSeconds int64
	MongoURL               string
	Port                   string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		DBHost:                 getEnv("DATABASE_HOST", "http://localhost"),
		DBUser:                 getEnv("DATABASE_USER", "root"),
		DBPassword:             getEnv("DATABASE_PASSWORD", "sabiq1234"),
		DBName:                 getEnv("DATABASE_NAME", "users"),
		JWTSecret:              getEnv("JWT_SECRET", "not-so-secret-now-is-it?"),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXPIRATION_IN_SECONDS", 3600*24*7),
		PETDB:                  getEnv("PETSTORE_DATABASE", "p"),
		MongoURL:               getEnv("MOGODB_URL", "password"),
		Port:                   getEnv("PORT", "8080"),
	}
}

// Gets the env by key or fallbacks
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}
