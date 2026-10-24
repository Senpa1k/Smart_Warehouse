// Package service defines business contracts and composes their implementations.
package service

import (
	"io"

	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/models"
	"github.com/Senpa1k/Smart_Warehouse/internal/repository"
	"github.com/Senpa1k/Smart_Warehouse/internal/service/services"
	"github.com/gorilla/websocket"
)

var made = make(chan any, 100) // до 100 запросов одновременно +- на время обработки запроса

// Authorization defines user registration and token operations.
type Authorization interface {
	CreateUser(models.Users) (uint, error)
	GetUser(string, string) (string, *models.Users, error)
	ParseToken(string) (uint, error)
}

// WebsocketDashBoard streams robot and inventory events to clients.
type WebsocketDashBoard interface {
	RunStream(*websocket.Conn)
}

// Inventory defines inventory import, export, and history operations.
type Inventory interface {
	ImportCSV(csvData io.Reader) (*entities.ImportResult, error)
	ExportExcel(productIDs []string) ([]byte, error)
	GetHistory(from, to, zone, status string, limit, offset int) (*entities.HistoryResponse, error)
}

// DashBoard loads the current warehouse dashboard snapshot.
type DashBoard interface {
	GetDashInfo(*entities.DashInfo) error
}

// Robot accepts scans and verifies registered robot IDs.
type Robot interface {
	AddData(entities.RobotsData) error
	CheckID(string) bool
}

// AI generates inventory forecasts.
type AI interface {
	Predict(entities.AIRequest) (*entities.AIResponse, error)
}

// Service groups business operations used by HTTP handlers.
type Service struct {
	Robot
	Inventory
	Authorization
	WebsocketDashBoard
	DashBoard
	AI
	Redis repository.Redis
}

// NewService wires business services to their repositories.
func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization:      services.NewAuthService(repos.Authorization),
		Robot:              services.NewRobotService(repos.Robot, made, repos.Redis),
		WebsocketDashBoard: services.NewWebsocketDashBoard(repos.WebsocketDashBoard, made),
		Inventory:          services.NewInventoryService(repos.Inventory, repos.Redis),
		DashBoard:          services.NewDashService(repos.DashBoard, repos.Redis),
		AI:                 services.NewAIService(repos.AI, made, repos.Redis),
		Redis:              repos.Redis,
	}
}
