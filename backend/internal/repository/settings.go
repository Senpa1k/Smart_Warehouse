package repository

import (
	"fmt"

	"github.com/Senpa1k/Smart_Warehouse/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// инициализация базы данных
func InitBD() (*gorm.DB, error) {
	dns, err := config.Get("DATABASE_URL")
	if err != nil {
		return nil, fmt.Errorf("faild conn %w", err)
	}
	my_db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("faild open %w", err)
	}

	return my_db, nil
}
