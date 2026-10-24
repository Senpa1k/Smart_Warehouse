package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/repository"
	"github.com/sirupsen/logrus"
)

// RobotService stores robot scans and publishes their updates.
type RobotService struct {
	repo  repository.Robot
	made  chan<- any
	redis repository.Redis
}

// NewRobotService constructs the robot service and its event channel.
func NewRobotService(repo repository.Robot, made chan<- any, redis repository.Redis) *RobotService {
	return &RobotService{
		repo:  repo,
		made:  made,
		redis: redis,
	}
}

// AddData persists a robot scan and updates its cached state.
func (r *RobotService) AddData(data entities.RobotsData) error {
	// Проверка валидности данных
	if !r.repo.CheckID(data.RobotID) {
		return fmt.Errorf("invalid robot id: %s", data.RobotID)
	}
	if err := r.repo.AddData(data); err != nil {
		return err
	}

	if r.redis != nil {
		if err := r.redis.SetRobotOnline(data.RobotID); err != nil {
			logrus.WithError(err).Warn("update robot online status")
		}
		if err := r.redis.SetRobotBattery(data.RobotID, data.BatteryLevel, 30*time.Second); err != nil {
			logrus.WithError(err).Warn("update robot battery")
		}
		if err := r.redis.SetRobotStatus(data.RobotID, "active", 30*time.Second); err != nil {
			logrus.WithError(err).Warn("update robot status")
		}
		logrus.Infof("Robot %s status updated in Redis", data.RobotID)
	}

	if r.redis != nil {
		if err := r.redis.Delete("dashboard:current"); err != nil {
			logrus.WithError(err).Warn("invalidate dashboard cache")
		}
	}

	// Публикуем событие в Redis
	if r.redis != nil {
		event := map[string]any{
			"type":      "robot_data",
			"robot_id":  data.RobotID,
			"battery":   data.BatteryLevel,
			"status":    "active",
			"online":    true,
			"timestamp": time.Now().Format(time.RFC3339),
		}

		eventJSON, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("marshal robot update: %w", err)
		}
		if err := r.redis.Publish("robot_updates", string(eventJSON)); err != nil {
			logrus.WithError(err).Warn("publish robot update")
		}
	}

	r.made <- data
	return nil
}

// CheckID reports whether a robot ID is registered.
func (r *RobotService) CheckID(robotID string) bool {
	return r.repo.CheckID(robotID)
}
