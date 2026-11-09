package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"task-management-system/internal/models"
)

// Connect opens the PostgreSQL connection and runs the auto-migrations.
func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Warn level: do not log every SQL statement (and its parameters).
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	if err := Migrate(db); err != nil {
		return nil, err
	}
	return db, nil
}

// Migrate creates or updates the schema.
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&models.User{},
		&models.TaskTag{},
		&models.Task{},
		&models.TaskTaskTag{},
	); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}
	return nil
}
