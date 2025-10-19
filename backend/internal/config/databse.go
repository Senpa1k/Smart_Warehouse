package config

import (
	"fmt"
	"log"
	"smart_warehouse/backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=postgres dbname=smart_warehouse port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Database connection established")

	// Пробуем миграции
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
