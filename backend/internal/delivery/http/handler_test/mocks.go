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

func (m *MockRobotService) AddData(data entities.RobotsData) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockRobotService) CheckId(id string) bool {
	args := m.Called(id)
	return args.Bool(0)
}

// MockDashboardService мок сервиса дашборда
type MockDashboardService struct {
	mock.Mock
}

func (m *MockDashboardService) GetDashInfo(dash *entities.DashInfo) error {
	args := m.Called(dash)
	return args.Error(0)
}

// MockWebsocketDashboardService мок вебсокет сервиса
type MockWebsocketDashboardService struct {
	mock.Mock
}

func (m *MockWebsocketDashboardService) RunStream(conn *websocket.Conn) {
	m.Called(conn)
}

// MockAIService мок AI сервиса
type MockAIService struct {
	mock.Mock
}

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

func (m *MockInventoryService) ExportExcel(ids []string) ([]byte, error) {
	args := m.Called(ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockInventoryService) ImportCSV(csvData io.Reader) (*entities.ImportResult, error) {
	args := m.Called(csvData)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ImportResult), args.Error(1)
}

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

func (m *MockRedisService) IsRobotOnline(robotID string) (bool, error) {
	args := m.Called(robotID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRedisService) GetRobotBattery(robotID string) (int, error) {
	args := m.Called(robotID)
	return args.Int(0), args.Error(1)
}

func (m *MockRedisService) GetRobotStatus(robotID string) (string, error) {
	args := m.Called(robotID)
	return args.String(0), args.Error(1)
}

func (m *MockRedisService) CheckRateLimit(key string, limit int, window time.Duration) (bool, error) {
	args := m.Called(key, limit, window)
	return args.Bool(0), args.Error(1)
}

func (m *MockRedisService) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockRedisService) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *MockRedisService) Set(key string, value interface{}, expiration time.Duration) error {
	args := m.Called(key, value, expiration)
	return args.Error(0)
}

// MockAuthService мок сервиса авторизации
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) CreateUser(user models.Users) (uint, error) {
	args := m.Called(user)
	return uint(args.Int(0)), args.Error(1)
}

func (m *MockAuthService) GetUser(email, password string) (string, *models.Users, error) {
	args := m.Called(email, password)
	var userPtr *models.Users
	if args.Get(1) != nil {
		user := args.Get(1).(models.Users)
		userPtr = &user
	}
	return args.String(0), userPtr, args.Error(2)
}

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

func (m *MockRedisService) Exists(key string) (bool, error) {
	args := m.Called(key)
	return args.Bool(0), args.Error(1)
}

func (m *MockRedisService) Publish(channel string, message interface{}) error {
	args := m.Called(channel, message)
	return args.Error(0)
}

func (m *MockRedisService) Subscribe(channel string) *redis.PubSub {
	args := m.Called(channel)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*redis.PubSub)
}

func (m *MockRedisService) SetRobotStatus(robotID, status string, expiration time.Duration) error {
	args := m.Called(robotID, status, expiration)
	return args.Error(0)
}

func (m *MockRedisService) SetRobotBattery(robotID string, batteryLevel int, expiration time.Duration) error {
	args := m.Called(robotID, batteryLevel, expiration)
	return args.Error(0)
}

func (m *MockRedisService) SetRobotOnline(robotID string) error {
	args := m.Called(robotID)
	return args.Error(0)
}
