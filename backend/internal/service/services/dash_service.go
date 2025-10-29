package services

import (
	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/repository"
)

type DashSevice struct {
	repo repository.DashBoard
}

func NewDashService(repo repository.DashBoard) *DashSevice {
	return &DashSevice{repo: repo}
}

func (d *DashSevice) GetDashInfo(dash *entities.DashInfo) error {
	return d.repo.GetDashInfo(dash)
}
