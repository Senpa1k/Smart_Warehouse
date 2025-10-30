package handler

import (
	"io"

	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/models"
	"github.com/Senpa1k/Smart_Warehouse/internal/service"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/mock"
)

// ==================== Моки для всех сервисов ====================

type MockService struct {
	mock.Mock
}

func (m *MockService) CreateUser(user models.Users) (uint, error) {
	args := m.Called(user)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockService) GetUser(email, password string) (string, *models.Users, error) {
	args := m.Called(email, password)
	return args.String(0), args.Get(1).(*models.Users), args.Error(2)
}

func (m *MockService) ParseToken(token string) (uint, error) {
	args := m.Called(token)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockService) RunStream(conn *websocket.Conn) {
	m.Called(conn)
}

func (m *MockService) ImportCSV(csvData io.Reader) (*entities.ImportResult, error) {
	args := m.Called(csvData)
	return args.Get(0).(*entities.ImportResult), args.Error(1)
}

func (m *MockService) ExportExcel(productIDs []string) ([]byte, error) {
	args := m.Called(productIDs)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockService) GetHistory(from, to, zone, status string, limit, offset int) (*entities.HistoryResponse, error) {
	args := m.Called(from, to, zone, status, limit, offset)
	return args.Get(0).(*entities.HistoryResponse), args.Error(1)
}

func (m *MockService) GetDashInfo(dash *entities.DashInfo) error {
	args := m.Called(dash)
	return args.Error(0)
}

func (m *MockService) AddData(data entities.RobotsData) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockService) CheckId(robotId string) bool {
	args := m.Called(robotId)
	return args.Bool(0)
}

func (m *MockService) Predict(request entities.AIRequest) (*entities.AIResponse, error) {
	args := m.Called(request)
	return args.Get(0).(*entities.AIResponse), args.Error(1)
}

// ==================== Адаптер ====================

type TestServiceWrapper struct {
	*MockService
}

func CreateTestService(mock *MockService) *service.Service {
	wrapper := &TestServiceWrapper{mock}

	return &service.Service{
		Robot:              wrapper,
		Inventory:          wrapper,
		Authorization:      wrapper,
		WebsocketDashBoard: wrapper,
		DashBoard:          wrapper,
		AI:                 wrapper,
	}
}

func NewTestHandler() (*Handler, *MockService) {
	mockService := new(MockService)
	testService := CreateTestService(mockService)
	handler := &Handler{services: testService}
	return handler, mockService
}
