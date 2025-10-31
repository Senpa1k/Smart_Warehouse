package services

import (
	"fmt"

	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/repository"
	"github.com/sirupsen/logrus"
)

type RobotService struct {
	repo  repository.Robot
	made  chan<- interface{}
	redis repository.Redis
}

func NewRobotService(repo repository.Robot, made chan<- interface{}, redis repository.Redis) *RobotService {
	return &RobotService{
		repo:  repo,
		made:  made,
		redis: redis,
	}
}

func (r *RobotService) AddData(data entities.RobotsData) error {
	// Проверка валидности данных
	if !r.repo.CheckId(data.RobotId) {
		return fmt.Errorf("invalid robot id: %s", data.RobotId)
	}

	// 1. Добавляем данные в БД
	if err := r.repo.AddData(data); err != nil {
		return err
	}

	// 2. Инвалидируем кеш dashboard при новых данных
	if r.redis != nil {
		// Удаляем кеш dashboard
		r.redis.Delete("dashboard:current")

		// Также можно инвалидировать кеш по зоне
		if data.Location.Zone != "" {
			cacheKey := fmt.Sprintf("dashboard:zone:%s", data.Location.Zone)
			r.redis.Delete(cacheKey)
		}

		logrus.Info("Dashboard cache invalidated due to new robot data")
	}

	// 3. Отправляем в канал для WebSocket
	r.made <- data

	return nil
}

func (r *RobotService) CheckId(robotID string) bool {
	return r.repo.CheckId(robotID)
}
