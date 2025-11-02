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

// добавление данных о сканировании
func (r *RobotService) AddData(data entities.RobotsData) error {
	// Проверка валидности данных
	if !r.repo.CheckId(data.RobotId) {
		return fmt.Errorf("invalid robot id: %s", data.RobotId)
	}
	if err := r.repo.AddData(data); err != nil {
		return err
	}

	if r.redis != nil {
		r.redis.SetRobotOnline(data.RobotId)
		r.redis.SetRobotBattery(data.RobotId, data.BatteryLevel, 30*time.Second)
		r.redis.SetRobotStatus(data.RobotId, "active", 30*time.Second)
		logrus.Infof("Robot %s status updated in Redis", data.RobotId)
	}

	if r.redis != nil {
		r.redis.Delete("dashboard:current")
	}

	//Публикуем событие в Redis
	if r.redis != nil {
		event := map[string]interface{}{
			"type":      "robot_data",
			"robot_id":  data.RobotId,
			"battery":   data.BatteryLevel,
			"status":    "active",
			"online":    true,
			"timestamp": time.Now().Format(time.RFC3339),
		}

		eventJSON, _ := json.Marshal(event)
		r.redis.Publish("robot_updates", string(eventJSON))
	}

	r.made <- data
	return nil
}

func (r *RobotService) CheckId(robotID string) bool {
	return r.repo.CheckId(robotID)
}
