package database

import (
	"fmt"
	"log"
	"time"

	"github.com/bhaskar/todo-api/internal/config"
	"github.com/bhaskar/todo-api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB holds the database connection
var DB *gorm.DB

// Connect establishes a database connection
func Connect(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	var dialector gorm.Dialector

	// Use SQLite for development/testing, PostgreSQL for production
	if cfg.Host == "sqlite" {
		dialector = sqlite.Open(cfg.DBName + ".db")
		log.Println("📦 Using SQLite database")
	} else {
		dialector = postgres.Open(cfg.DSN())
		log.Println("🐘 Connecting to PostgreSQL...")
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Store globally
	DB = db

	log.Println("✅ Database connected successfully")
	return db, nil
}

// Migrate runs auto-migrations for all models
func Migrate(db *gorm.DB) error {
	log.Println("🔄 Running database migrations...")
	
	err := db.AutoMigrate(
		&models.User{},
		&models.Todo{},
	)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("✅ Database migration completed")
	return nil
}

// Close closes the database connection
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
