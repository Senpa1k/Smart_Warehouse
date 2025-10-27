package service

import (
	"io"

	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/models"
	"github.com/Senpa1k/Smart_Warehouse/internal/repository"
)

type Authorization interface {
	CreateUser(models.Users) (uint, error)
	GetUser(string, string) (string, *models.Users, error)
	ParseToken(string) (uint, error)
}

type DashBoard interface {
}

type Inventory interface {
	ImportCSV(csvData io.Reader) (*ImportResult, error)
	ExportExcel(productIDs []string) ([]byte, error)
	GetHistory(from, to, zone, status string, limit, offset int) (*HistoryResponse, error)
}

type Robot interface {
	AddData(entities.RobotsData) error
}

type Service struct {
	Robot
	Inventory
	Authorization
	DashBoard
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Robot:         NewRobotService(repos.Robot),
		Inventory:     NewInventoryService(repos.Inventory),
	}
}
