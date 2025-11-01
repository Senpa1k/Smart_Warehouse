package services

import (
	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/repository"
)

type RobotService struct {
	repo repository.Robot
	made chan<- interface{}
}

func NewRobotService(repo repository.Robot, made chan<- interface{}) *RobotService {
	return &RobotService{repo: repo, made: made}
}

// добавление данных о сканировании
func (r *RobotService) AddData(data entities.RobotsData) error {
	//проверка валидности данных

	if err := r.repo.AddData(data); err != nil {
		return err
	}

	r.made <- data

	return nil
}

func (r *RobotService) CheckId(robotID string) bool {
	return r.repo.CheckId(robotID)
}
