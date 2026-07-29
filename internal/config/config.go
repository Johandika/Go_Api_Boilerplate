package config

import (
	"bufio"
	"os"
	"strings"
)

type Config struct {
	AppPort                  string
	AppEnv                   string
	DatabaseURL              string
	APIKey                   string
	JWTSecret                string
	JWTAccessTokenExpiresIn  string
	JWTRefreshTokenExpiresIn string
	CORSAllowedOrigins       string
}

func Load() Config {
	loadDotEnv(".env")
	return Config{
		AppPort:                  getEnv("APP_PORT", "8080"),
		AppEnv:                   getEnv("APP_ENV", "development"),
		DatabaseURL:              getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/go_api_boilerplate?sslmode=disable"),
		APIKey:                   getEnv("API_KEY", "change-this-api-key"),
		JWTSecret:                getEnv("JWT_SECRET", "change-this-jwt-secret"),
		JWTAccessTokenExpiresIn:  getEnv("JWT_ACCESS_TOKEN_EXPIRES_IN", "15m"),
		JWTRefreshTokenExpiresIn: getEnv("JWT_REFRESH_TOKEN_EXPIRES_IN", "168h"),
		CORSAllowedOrigins:       getEnv("CORS_ALLOWED_ORIGINS", "*"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func loadDotEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), "\"'")
		if os.Getenv(key) == "" {
			_ = os.Setenv(key, value)
		}
	}
}
