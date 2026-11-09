package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds runtime configuration loaded from the environment / .env file.
type Config struct {
	DBURI          string
	JWTSecret      string
	GeminiKey      string
	GeminiModel    string
	Host           string
	Port           string
	AllowedOrigins []string
}

const minJWTSecretLength = 32

// LoadConfig reads configuration from the environment. Required values have
// no fallbacks: a missing DB_URI or JWT_SECRET is reported as an error.
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg := &Config{
		DBURI:          strings.TrimSpace(os.Getenv("DB_URI")),
		JWTSecret:      strings.TrimSpace(os.Getenv("JWT_SECRET")),
		GeminiKey:      strings.TrimSpace(strings.Trim(os.Getenv("GEMINI_KEY"), "\"")),
		GeminiModel:    getEnvWithDefault("GEMINI_MODEL", "gemini-2.5-flash"),
		Host:           getEnvWithDefault("HOST", "127.0.0.1"),
		Port:           getEnvWithDefault("PORT", "8080"),
		AllowedOrigins: splitList(getEnvWithDefault("CORS_ORIGINS", "http://localhost:3000")),
	}

	var missing []string
	if cfg.DBURI == "" {
		missing = append(missing, "DB_URI")
	}
	if cfg.JWTSecret == "" {
		missing = append(missing, "JWT_SECRET")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variable(s): %s (see .env.example)", strings.Join(missing, ", "))
	}
	if len(cfg.JWTSecret) < minJWTSecretLength {
		return nil, errors.New("JWT_SECRET must be at least 32 characters (try: openssl rand -hex 32)")
	}
	if cfg.GeminiKey == "" {
		log.Println("GEMINI_KEY is not set: AI suggestions are disabled")
	}

	return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return defaultValue
}

func splitList(value string) []string {
	var out []string
	for _, item := range strings.Split(value, ",") {
		if item = strings.TrimSpace(item); item != "" {
			out = append(out, strings.TrimRight(item, "/"))
		}
	}
	return out
}
