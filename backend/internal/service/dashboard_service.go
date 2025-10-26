package service

import (
	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/repository"
)

type DashBoardService struct {
	repo repository.DashBoard
}

func NewDashBoardService(repo repository.DashBoard) *DashBoardService {
	return &DashBoardService{
		repo: repo,
	}
}

func (d *DashBoardService) GetInfo() (*entities.DashInfo, error) {
	return d.repo.GetInfo()
}
