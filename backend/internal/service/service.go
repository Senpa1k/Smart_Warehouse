package service

import (
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

type History interface {
}

type Inventory interface {
}

type Robot interface {
	AddData(entities.RobotsData) error
	CheckId(string) bool
}

type Service struct {
	Robot
	Inventory
	History
	Authorization
	WebsocketDashBoard
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization:      NewAuthService(repos.Authorization),
		Robot:              NewRobotService(repos.Robot, made),
		WebsocketDashBoard: NewWebsocketDashBoard(repos.WebsocketDashBoard, made),
	}
}
