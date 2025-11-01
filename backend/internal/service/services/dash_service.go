package services

import (
	"encoding/json"
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/repository"
	"github.com/sirupsen/logrus"
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
func (d *DashSevice) GetDashInfo(dash *entities.DashInfo) error {
	return d.repo.GetDashInfo(dash)
}
