package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"time"

	"task-management-system/internal/auth"
	"task-management-system/internal/config"
	"task-management-system/internal/database"
	"task-management-system/internal/routes"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}
	auth.Init(cfg.JWTSecret)

	db, err := database.Connect(cfg.DBURI)
	if err != nil {
		log.Fatalf("Database error: %v", err)
	}

	router := routes.SetupRouter(db, cfg)

	addr := net.JoinHostPort(cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	log.Printf("Server listening on http://%s (CORS origins: %v)", addr, cfg.AllowedOrigins)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
