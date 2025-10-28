package service

import (
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

type History interface {
}

type Inventory interface {
}

type Robot interface {
	AddData(entities.RobotsData) error
}

type AI interface {
	Predict(entities.AIRequest) (*entities.AIResponse, error)
}

type Service struct {
	Robot
	Inventory
	History
	Authorization
	DashBoard
	AI
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Robot:         NewRobotService(repos.Robot),
		AI: 			NewAIService(repos.AI),
	}
}
