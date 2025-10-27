package repository

import (
	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/models"
	"gorm.io/gorm"
)

type Authorization interface {
	CreateUser(models.Users) (uint, error)
	GetUser(string, string) (*models.Users, error)
}

type DashBoard interface {
}

type Inventory interface {
	ImportInventoryHistories(histories []models.InventoryHistory) error
	GetInventoryHistoryByProductIDs(productIDs []string) ([]models.InventoryHistory, error)
	GetProductByID(productID string) error
	CreateProduct(product *models.Products) error
	UpdateProduct(product *models.Products) error
	GetHistory(from, to, zone, status string, limit, offset int) ([]models.InventoryHistory, int64, error)
}

type Robot interface {
	AddData(entities.RobotsData) error
}

type Repository struct {
	Robot
	Inventory
	Authorization
	DashBoard
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Robot:         NewRobotPostgres(db),
		Inventory:     NewInventoryRepo(db),
	}
}
