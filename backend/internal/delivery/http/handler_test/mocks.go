package test_handler

import (
	"io"
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/models"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
)

// MockRobotService мок сервиса роботов
type MockRobotService struct {
	mock.Mock
}

// AddData records a robot scan submission.
func (m *MockRobotService) AddData(data entities.RobotsData) error {
	args := m.Called(data)
	return args.Error(0)
}

// CheckID returns the configured robot registration result.
func (m *MockRobotService) CheckID(id string) bool {
	args := m.Called(id)
	return args.Bool(0)
}

// MockDashboardService мок сервиса дашборда
type MockDashboardService struct {
	mock.Mock
}

// GetDashInfo fills a dashboard response from the configured mock call.
func (m *MockDashboardService) GetDashInfo(dash *entities.DashInfo) error {
	args := m.Called(dash)
	return args.Error(0)
}

// MockWebsocketDashboardService мок вебсокет сервиса
type MockWebsocketDashboardService struct {
	mock.Mock
}

// RunStream records a WebSocket stream invocation.
func (m *MockWebsocketDashboardService) RunStream(conn *websocket.Conn) {
	m.Called(conn)
}

// MockAIService мок AI сервиса
type MockAIService struct {
	mock.Mock
}

// Predict returns the configured forecast or error.
func (m *MockAIService) Predict(request entities.AIRequest) (*entities.AIResponse, error) {
	args := m.Called(request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.AIResponse), args.Error(1)
}

// MockInventoryService мок сервиса инвентаря
type MockInventoryService struct {
	mock.Mock
}

// ExportExcel returns the configured Excel payload or error.
func (m *MockInventoryService) ExportExcel(ids []string) ([]byte, error) {
	args := m.Called(ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

// ImportCSV returns the configured CSV import result.
func (m *MockInventoryService) ImportCSV(csvData io.Reader) (*entities.ImportResult, error) {
	args := m.Called(csvData)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ImportResult), args.Error(1)
}

// GetHistory returns the configured inventory history.
func (m *MockInventoryService) GetHistory(from, to, zone, status string, limit, offset int) (*entities.HistoryResponse, error) {
	args := m.Called(from, to, zone, status, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.HistoryResponse), args.Error(1)
}

// MockRedisService мок Redis сервиса
type MockRedisService struct {
	mock.Mock
}

// IsRobotOnline returns the configured online state.
func (m *MockRedisService) IsRobotOnline(robotID string) (bool, error) {
	args := m.Called(robotID)
	return args.Bool(0), args.Error(1)
}

// GetRobotBattery returns the configured battery level.
func (m *MockRedisService) GetRobotBattery(robotID string) (int, error) {
	args := m.Called(robotID)
	return args.Int(0), args.Error(1)
}

// GetRobotStatus returns the configured robot status.
func (m *MockRedisService) GetRobotStatus(robotID string) (string, error) {
	args := m.Called(robotID)
	return args.String(0), args.Error(1)
}

// CheckRateLimit returns the configured rate-limit decision.
func (m *MockRedisService) CheckRateLimit(key string, limit int, window time.Duration) (bool, error) {
	args := m.Called(key, limit, window)
	return args.Bool(0), args.Error(1)
}

// Delete records a cache deletion.
func (m *MockRedisService) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

// Get returns a configured cache value.
func (m *MockRedisService) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

// Set records a cache write.
func (m *MockRedisService) Set(key string, value any, expiration time.Duration) error {
	args := m.Called(key, value, expiration)
	return args.Error(0)
}

// MockAuthService мок сервиса авторизации
type MockAuthService struct {
	mock.Mock
}

// CreateUser returns the configured new user ID or error.
func (m *MockAuthService) CreateUser(user models.Users) (uint, error) {
	args := m.Called(user)
	return uint(args.Int(0)), args.Error(1)
}

// GetUser returns the configured user account and token.
func (m *MockAuthService) GetUser(email, password string) (string, *models.Users, error) {
	args := m.Called(email, password)
	var userPtr *models.Users
	if args.Get(1) != nil {
		user := args.Get(1).(models.Users)
		userPtr = &user
	}
	return args.String(0), userPtr, args.Error(2)
}

// ParseToken returns the configured token validation result.
func (m *MockAuthService) ParseToken(token string) (uint, error) {
	args := m.Called(token)
	return uint(args.Int(0)), args.Error(1)
}

// MockServices мок всех сервисов
type MockServices struct {
	Robot              *MockRobotService
	DashBoard          *MockDashboardService
	WebsocketDashBoard *MockWebsocketDashboardService
	AI                 *MockAIService
	Inventory          *MockInventoryService
	Redis              *MockRedisService
	Authorization      *MockAuthService
}

// NewMockServices создает новый экземпляр мок сервисов
func NewMockServices() *MockServices {
	return &MockServices{
		Robot:              new(MockRobotService),
		DashBoard:          new(MockDashboardService),
		WebsocketDashBoard: new(MockWebsocketDashboardService),
		AI:                 new(MockAIService),
		Inventory:          new(MockInventoryService),
		Redis:              new(MockRedisService),
		Authorization:      new(MockAuthService),
	}
}

// Exists returns the configured cache presence result.
func (m *MockRedisService) Exists(key string) (bool, error) {
	args := m.Called(key)
	return args.Bool(0), args.Error(1)
}

// Publish records a Redis publication.
func (m *MockRedisService) Publish(channel string, message any) error {
	args := m.Called(channel, message)
	return args.Error(0)
}

// Subscribe returns the configured Redis subscription.
func (m *MockRedisService) Subscribe(channel string) *redis.PubSub {
	args := m.Called(channel)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*redis.PubSub)
}

// SetRobotStatus records a robot status update.
func (m *MockRedisService) SetRobotStatus(robotID, status string, expiration time.Duration) error {
	args := m.Called(robotID, status, expiration)
	return args.Error(0)
}

// SetRobotBattery records a robot battery update.
func (m *MockRedisService) SetRobotBattery(robotID string, batteryLevel int, expiration time.Duration) error {
	args := m.Called(robotID, batteryLevel, expiration)
	return args.Error(0)
}

// SetRobotOnline records a robot heartbeat update.
func (m *MockRedisService) SetRobotOnline(robotID string) error {
	args := m.Called(robotID)
	return args.Error(0)
}
