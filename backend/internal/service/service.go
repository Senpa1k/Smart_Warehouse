package service

import (
	"github.com/Senpa1k/Smart_Warehouse/internal/models"
	"github.com/Senpa1k/Smart_Warehouse/internal/repository"
)

type Authorization interface {
	CreateUser(models.Users) (uint, error)
	GetUser(string, string) (string, *models.Users, error)
}

type DashBoard interface {
}

type History interface {
}

type Inventory interface {
}

type Robot interface {
}

type Service struct {
	Robot
	Inventory
	History
	Authorization
	DashBoard
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
