package repository

import (
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/models"
	"gorm.io/gorm"
)

type Authorization interface {
	CreateUser(models.Users) (uint, error)
	GetUser(string, string) (*models.Users, error)
}

type WebsocketDashBoard interface {
	InventoryAlertScanned(*entities.InventoryAlert, time.Time, string) error
}

type History interface {
}

type Inventory interface {
}

type Robot interface {
	AddData(entities.RobotsData) error
	CheckId(string) bool
}

type Repository struct {
	Robot
	Inventory
	History
	Authorization
	WebsocketDashBoard
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Authorization:      NewAuthPostgres(db),
		Robot:              NewRobotPostges(db),
		WebsocketDashBoard: NewWebsocketDashBoardPostgres(db),
	}
}
