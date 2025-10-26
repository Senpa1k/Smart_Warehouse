package repository

import (
	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"gorm.io/gorm"
)

type DashBoardPosgres struct {
	db *gorm.DB
}

func NewDashBoardPostges(db *gorm.DB) *DashBoardPosgres {
	return &DashBoardPosgres{
		db: db,
	}
}

func (d *DashBoardPosgres) GetInfo() (*entities.DashInfo, error) {
	return nil, nil
}
