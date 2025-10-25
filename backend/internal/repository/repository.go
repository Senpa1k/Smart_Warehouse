package repository

import "gorm.io/gorm"

type Authorization interface {
}

type DashBoard interface {
}

type History interface {
}

type Inventory interface {
}

type Robot interface {
}

type Repository struct {
	Robot
	Inventory
	History
	Authorization
	DashBoard
}

func NewRepository(*gorm.DB) *Repository {
	return &Repository{}
}
