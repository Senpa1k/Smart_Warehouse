package repository

import (
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/models"
	"github.com/Senpa1k/Smart_Warehouse/internal/repository/postgres"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Authorization defines persistence operations for user accounts.
type Authorization interface {
	CreateUser(models.Users) (uint, error)
	GetUser(string, string) (*models.Users, error)
}

// WebsocketDashBoard loads inventory alerts for dashboard streams.
type WebsocketDashBoard interface {
	InventoryAlertScanned(*entities.InventoryAlert, time.Time, string) error
	InventoryAlertPredict(*entities.InventoryAlert, entities.Predictions) error
}

// Inventory defines inventory persistence and history queries.
type Inventory interface {
	ImportInventoryHistories(histories []models.InventoryHistory) error
	GetInventoryHistoryByProductIDs(productIDs []string) ([]models.InventoryHistory, error)
	GetInventoryHistoryByScanIDs(scanIDs []string) ([]models.InventoryHistory, error)
	GetProductByID(productID string) error
	CreateProduct(product *models.Products) error
	UpdateProduct(product *models.Products) error
	GetHistory(from, to, zone, status string, limit, offset int) ([]models.InventoryHistory, int64, error)
}

// DashBoard loads data for the current dashboard snapshot.
type DashBoard interface {
	GetDashInfo(*entities.DashInfo) error
}

// Robot stores robot scans and checks registered IDs.
type Robot interface {
	AddData(entities.RobotsData) error
	CheckID(string) bool
}

// AI loads AI inputs and stores generated forecasts.
type AI interface {
	AIRequest(entities.AIRequest) ([]models.InventoryHistory, error)
	AIResponse(entities.AIResponse) error
}

// Redis интерфейс
type Redis interface {
	Set(key string, value any, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
	Exists(key string) (bool, error)
	Publish(channel string, message any) error
	Subscribe(channel string) *redis.PubSub

	SetRobotStatus(robotID, status string, expiration time.Duration) error
	GetRobotStatus(robotID string) (string, error)
	SetRobotBattery(robotID string, batteryLevel int, expiration time.Duration) error
	GetRobotBattery(robotID string) (int, error)
	SetRobotOnline(robotID string) error
	IsRobotOnline(robotID string) (bool, error)

	CheckRateLimit(key string, limit int, window time.Duration) (bool, error)
}

// Repository groups the persistence interfaces used by services.
type Repository struct {
	Robot
	Inventory
	Authorization
	WebsocketDashBoard
	DashBoard
	AI
	Redis Redis
}

// NewRepository wires PostgreSQL repositories and an optional Redis client.
func NewRepository(db *gorm.DB, redisClient Redis) *Repository {
	return &Repository{
		Authorization:      postgres.NewAuthPostgres(db),
		Robot:              postgres.NewRobotPostgres(db),
		WebsocketDashBoard: postgres.NewWebsocketDashBoardPostgres(db),
		Inventory:          postgres.NewInventoryRepo(db),
		DashBoard:          postgres.NewDashPostgres(db),
		AI:                 postgres.NewAIPostgres(db),
		Redis:              redisClient,
	}
}

// WithRedis runs fn when a Redis client is configured.
func (r *Repository) WithRedis(fn func(redis Redis) error) error {
	if r.Redis == nil {
		return nil
	}
	return fn(r.Redis)
}
