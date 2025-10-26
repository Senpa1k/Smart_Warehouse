package service

import (
	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/repository"
)

type RobotService struct {
	repo repository.Robot
}

func NewRobotService(repo repository.Robot) *RobotService {
	return &RobotService{repo: repo}
}

func (r *RobotService) AddData(data entities.RobotsData) error {
	//проверка валидности данных
	return r.repo.AddData(data)
}
