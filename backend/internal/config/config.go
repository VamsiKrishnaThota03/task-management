package config

import (
    "log"
    "os"
    "strings"
    "github.com/joho/godotenv"
)

type Config struct {
    DBURI      string
    JWTSecret  string
    GeminiKey  string
}

func LoadConfig() *Config {
    if err := godotenv.Load(); err != nil {
        log.Println("Warning: No .env file found, using system environment variables")
    }

    cfg := &Config{
        DBURI:     getEnvWithDefault("DB_URI", ""),
        JWTSecret: getEnvWithDefault("JWT_SECRET", ""),
        GeminiKey: strings.TrimSpace(strings.Trim(getEnvWithDefault("GEMINI_KEY", ""), "\"")),
    }

    // Validate required configuration
    if cfg.GeminiKey == "" {
        log.Fatal("Error: GEMINI_KEY is required but not set")
    }

    log.Printf("Configuration loaded successfully")
    return cfg
}

func getEnvWithDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
