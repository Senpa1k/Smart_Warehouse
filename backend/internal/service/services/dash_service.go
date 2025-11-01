package services

import (
	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/repository"
)

type DashService struct {
	repo  repository.DashBoard
	redis repository.Redis
}

func NewDashService(repo repository.DashBoard, redis repository.Redis) *DashService {
	return &DashService{
		repo:  repo,
		redis: redis,
	}
}

// получение данных о карте
func (d *DashService) GetDashInfo(dash *entities.DashInfo) error {
	return d.repo.GetDashInfo(dash)
}
