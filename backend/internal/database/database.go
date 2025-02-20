package database

import (
    "log"
    "task-management-system/internal/config"
    "task-management-system/internal/models"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB(cfg *config.Config) *gorm.DB {
    db, err := gorm.Open(postgres.Open(cfg.DBURI), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Auto-migrate the schemas in correct order
    err = db.AutoMigrate(
        &models.User{},
        &models.TaskTag{},
        &models.Task{},
        &models.TaskTaskTag{}, // Add the join table to migrations
    )
    if err != nil {
        log.Fatal("Failed to run migrations:", err)
    }

    return db
}
