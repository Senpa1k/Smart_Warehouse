package services

import (
	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/repository"
)

// DashService loads the current warehouse dashboard snapshot.
type DashService struct {
	repo  repository.DashBoard
	redis repository.Redis
}

// NewDashService constructs the dashboard service.
func NewDashService(repo repository.DashBoard, redis repository.Redis) *DashService {
	return &DashService{
		repo:  repo,
		redis: redis,
	}
}

// GetDashInfo fills a dashboard response from the repository.
func (d *DashService) GetDashInfo(dash *entities.DashInfo) error {
	return d.repo.GetDashInfo(dash)
}
