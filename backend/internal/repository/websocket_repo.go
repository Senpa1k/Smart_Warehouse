package repository

import (
	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"gorm.io/gorm"
)

type WebsocketDashBoardPostgres struct {
	db *gorm.DB
}

func NewWebsocketDashBoardPostgres(db *gorm.DB) *WebsocketDashBoardPostgres {
	return &WebsocketDashBoardPostgres{db: db}
}

func (w *WebsocketDashBoardPostgres) UpdateRobot(entitie *entities.UpdateRobot) error {
	return nil
}

func (w *WebsocketDashBoardPostgres) InventoryAlert(entitie *entities.InventoryAlert) error {
	return nil
}
