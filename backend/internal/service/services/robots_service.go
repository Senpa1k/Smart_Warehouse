package services

import (
	"encoding/json"
	"fmt"
	"time"

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

	// 2. Инвалидируем кеш dashboard
	if r.redis != nil {
		r.redis.Delete("dashboard:current")
		logrus.Info("Dashboard cache invalidated due to new robot data")
	}

	// 3. ✅ НОВОЕ: Публикуем событие в Redis
	if r.redis != nil {
		event := map[string]interface{}{
			"type":      "robot_data",
			"data":      data,
			"timestamp": time.Now().Format(time.RFC3339),
		}

		eventJSON, _ := json.Marshal(event)
		r.redis.Publish("robot_updates", string(eventJSON))
		logrus.Infof("📤 Published robot data to Redis channel: %s", data.RobotId)
	}

	// 4. Отправляем в канал для WebSocket (старая логика)
	r.made <- data

	return nil
}

func (r *RobotService) CheckId(robotID string) bool {
	return r.repo.CheckId(robotID)
}
