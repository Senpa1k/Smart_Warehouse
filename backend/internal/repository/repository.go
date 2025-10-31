package repository

import (
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/models"
	"github.com/Senpa1k/Smart_Warehouse/internal/repository/postgres"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Authorization interface {
	CreateUser(models.Users) (uint, error)
	GetUser(string, string) (*models.Users, error)
}

type WebsocketDashBoard interface {
	InventoryAlertScanned(*entities.InventoryAlert, time.Time, string) error
	InventoryAlertPredict(*entities.InventoryAlert, entities.Predictions) error
}

type Inventory interface {
	ImportInventoryHistories(histories []models.InventoryHistory) error
	GetInventoryHistoryByProductIDs(productIDs []string) ([]models.InventoryHistory, error)
	GetInventoryHistoryByScanIDs(scanIDs []string) ([]models.InventoryHistory, error)
	GetProductByID(productID string) error
	CreateProduct(product *models.Products) error
	UpdateProduct(product *models.Products) error
	GetHistory(from, to, zone, status string, limit, offset int) ([]models.InventoryHistory, int64, error)
}

type DashBoard interface {
	GetDashInfo(*entities.DashInfo) error
}

type Robot interface {
	AddData(entities.RobotsData) error
	CheckId(string) bool
}

type AI interface {
	AIRequest(entities.AIRequest) (*[]models.Products, error)
	AIResponse(entities.AIResponse) error
}

// Redis интерфейс
type Redis interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
	Exists(key string) (bool, error)
	Publish(channel string, message interface{}) error
	Subscribe(channel string) *redis.PubSub

	SetRobotStatus(robotID, status string, expiration time.Duration) error
	GetRobotStatus(robotID string) (string, error)
	SetRobotBattery(robotID string, batteryLevel int, expiration time.Duration) error
	GetRobotBattery(robotID string) (int, error)
	SetRobotOnline(robotID string) error
	IsRobotOnline(robotID string) (bool, error)

	CheckRateLimit(key string, limit int, window time.Duration) (bool, error)
}

type Repository struct {
	Robot
	Inventory
	Authorization
	WebsocketDashBoard
	DashBoard
	AI
	Redis Redis // Добавляем Redis
}

// Меняем конструктор чтобы принимать Redis клиент
func NewRepository(db *gorm.DB, redisClient Redis) *Repository {
	return &Repository{
		Authorization:      postgres.NewAuthPostgres(db),
		Robot:              postgres.NewRobotPostgres(db),
		WebsocketDashBoard: postgres.NewWebsocketDashBoardPostgres(db),
		Inventory:          postgres.NewInventoryRepo(db),
		DashBoard:          postgres.NewDashPostgres(db),
		AI:                 postgres.NewAIPostgres(db),
		Redis:              redisClient, // Передаем Redis клиент
	}
}

// Helper метод для безопасной работы с Redis
func (r *Repository) WithRedis(fn func(redis Redis) error) error {
	if r.Redis == nil {
		// Redis не доступен, пропускаем операцию
		return nil
	}
	return fn(r.Redis)
}
