package config

import (
	"fmt"
	"log"
	"os"

	"github.com/Senpa1k/Smart_Warehouse/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getEnv(key string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return ""
}

func InitDB() (*gorm.DB, error) {
	dsn := getEnv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Database connection established")

	err = db.AutoMigrate(
		&models.Users{},
		&models.Products{},
		&models.Robots{},
	)
	if err != nil {
		return nil, fmt.Errorf("migration error: %w", err)
	}
	err = db.AutoMigrate(
		&models.InventoryHistory{},
		&models.AiPrediction{},
	)
	if err != nil {
		return nil, fmt.Errorf("migration error: %w", err)
	}

	log.Println("Auto migration completed successfully")
	return db, nil
}
