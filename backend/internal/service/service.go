package service

import (
	"io"

	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/models"
	"github.com/Senpa1k/Smart_Warehouse/internal/repository"
	"github.com/gorilla/websocket"
)

var made = make(chan interface{}, 100) // до 100 запросов одновременно +- на время обработки запроса

type Authorization interface {
	CreateUser(models.Users) (uint, error)
	GetUser(string, string) (string, *models.Users, error)
	ParseToken(string) (uint, error)
}

type WebsocketDashBoard interface {
	RunStream(*websocket.Conn)
}

type Inventory interface {
	ImportCSV(csvData io.Reader) (*ImportResult, error)
	ExportExcel(productIDs []string) ([]byte, error)
	GetHistory(from, to, zone, status string, limit, offset int) (*HistoryResponse, error)
}

type DashBoard interface {
	GetDashInfo(*entities.DashInfo) error
}

type Robot interface {
	AddData(entities.RobotsData) error
	CheckId(string) bool
}

type AI interface {
	Predict(entities.AIRequest) (*entities.AIResponse, error)
}

type Service struct {
	Robot
	Inventory
	Authorization
	WebsocketDashBoard
	DashBoard
	AI
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization:      NewAuthService(repos.Authorization),
		Robot:              NewRobotService(repos.Robot, made),
		WebsocketDashBoard: NewWebsocketDashBoard(repos.WebsocketDashBoard, made),
		Inventory:          NewInventoryService(repos.Inventory),
		DashBoard:          NewDashService(repos.DashBoard),
		AI:                 NewAIService(repos.AI),
	}
}
