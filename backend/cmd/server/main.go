package main

import (
    "log"
    "net/http"
    "task-management-system/internal/config"
    "task-management-system/internal/database"
    "task-management-system/internal/routes"
)

func main() {
    // Load configuration
    cfg := config.LoadConfig()
    
    // Connect to database
    db := database.ConnectDB(cfg)
    
    // Setup router with both db and config
    router := routes.SetupRouter(db, cfg)

    log.Printf("Server configuration loaded:")
    log.Printf("DB_URI: %s", cfg.DBURI)
    log.Printf("JWT_SECRET: %s", cfg.JWTSecret)
    log.Printf("GEMINI_KEY: %s", cfg.GeminiKey)

    log.Printf("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}