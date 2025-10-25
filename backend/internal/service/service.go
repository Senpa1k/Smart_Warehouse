package service

import "github.com/Senpa1k/Smart_Warehouse/internal/repository"

type Authorization interface {
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
	return &Service{}
}
