package repository

import (
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/models"
	"github.com/Senpa1k/Smart_Warehouse/internal/repository/postgres"
	"gorm.io/gorm"
)

type Authorization interface {
	CreateUser(models.Users) (uint, error)
	GetUser(string, string) (*models.Users, error)
}

type WebsocketDashBoard interface {
	InventoryAlertScanned(*entities.InventoryAlert, time.Time, string) error
}

type Inventory interface {
	ImportInventoryHistories(histories []models.InventoryHistory) error
	GetInventoryHistoryByProductIDs(productIDs []string) ([]models.InventoryHistory, error)
	GetProductByID(productID string) error
	CreateProduct(product *models.Products) error
	UpdateProduct(product *models.Products) error
	GetHistory(from, to, zone, status string, limit, offset int) ([]models.InventoryHistory, int64, error)
}

type DashBoard interface {
	GetDashInfo(*entities.DashInfo) error
}

type Robot interface {
	AddData(entities.RobotsData) error
	CheckId(string) bool
}

type AI interface {
	AIRequest(entities.AIRequest) (*[]models.Products, error)
	AIResponse(entities.AIResponse) error
}

type Repository struct {
	Robot
	Inventory
	Authorization
	WebsocketDashBoard
	DashBoard
	AI
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Authorization:      postgres.NewAuthPostgres(db),
		Robot:              postgres.NewRobotPostgres(db),
		WebsocketDashBoard: postgres.NewWebsocketDashBoardPostgres(db),
		Inventory:          postgres.NewInventoryRepo(db),
		DashBoard:          postgres.NewDashPostgres(db),
		AI:                 postgres.NewAIPostgres(db),
	}
}
