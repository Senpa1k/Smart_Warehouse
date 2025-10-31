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

func (s *DashService) GetDashInfo(dash *entities.DashInfo) error {
	// Проверить кеш
	if s.redis != nil {
		if cached, err := s.redis.Get("dashboard:current"); err == nil {
			logrus.Info("Dashboard data served from cache")
			return json.Unmarshal([]byte(cached), dash)
		}
	}

	// Получить из БД
	err := s.repo.GetDashInfo(dash)

	// Сохранить в кеш на 5 секунд
	if err == nil && s.redis != nil {
		data, _ := json.Marshal(dash)
		s.redis.Set("dashboard:current", data, 5*time.Second)
		logrus.Info("Dashboard data cached")
	}

	return err
}
