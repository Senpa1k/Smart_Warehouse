// Package models maps warehouse records to database tables.
package models

import (
	"time"
)

// структыры таблиц бд для сопоставления с gorm

// Users maps a warehouse user account to the users table.
type Users struct {
	ID           uint      `gorm:"primaryKey;autoIncrement;type:serial" json:"id"`
	Email        string    `gorm:"unique;not null;type:varchar(255)" json:"email" binding:"required"`
	PasswordHash string    `gorm:"not null;type:varchar(255)" json:"password" binding:"required"`
	Name         string    `gorm:"not null;type:varchar(255)" json:"name"`
	Role         string    `gorm:"size:50;not null" json:"role"` // isv_valid   --'operator', 'admin', 'viewer'
	CreatedAt    time.Time `gorm:"type:timestamptz;default:now()" json:"created_at"`
}

// Products maps a catalog product and its stock thresholds.
type Products struct {
	ID           string `gorm:"primaryKey;type:varchar(50)" json:"id"`
	Name         string `gorm:"type:varchar(255);not null" json:"name"`
	Category     string `gorm:"type:varchar(100)" json:"category"`
	MinStock     int    `gorm:"default:10" json:"min_stock"`
	OptimalStock int    `gorm:"default:100" json:"optimal_stock"`
}

// Robots maps a robot and its most recent location and battery state.
type Robots struct {
	ID           string    `gorm:"primaryKey;type:varchar(50)" json:"id"`
	Status       string    `gorm:"size:50;default:active" json:"status"`
	BatteryLevel int       `gorm:"type:int" json:"battery_level"`
	LastUpdate   time.Time `gorm:"type:timestamptz" json:"last_update"`
	CurrentZone  string    `gorm:"size:10" json:"current_zone"`
	CurrentRow   int       `gorm:"type:int" json:"current_row"`
	CurrentShelf int       `gorm:"type:int" json:"current_shelf"`
}

// InventoryHistory maps an inventory scan and its related robot and product.
type InventoryHistory struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	RobotID     string    `gorm:"type:varchar(50);not null" json:"robot_id"`
	ProductID   string    `gorm:"type:varchar(50);not null" json:"product_id"`
	Quantity    int       `gorm:"type:integer;not null" json:"quantity"`
	Zone        string    `gorm:"size:10;not null" json:"zone"`
	RowNumber   int       `gorm:"type:integer" json:"row_number"`
	ShelfNumber int       `gorm:"type:integer" json:"shelf_number"`
	Status      string    `gorm:"size:50" json:"status"`
	ScannedAt   time.Time `gorm:"type:timestamptz;not null" json:"scanned_at"`
	CreatedAt   time.Time `gorm:"type:timestamptz;default:now()" json:"created_at"`

	// Связи
	Robot   Robots   `gorm:"foreignKey:RobotID;references:ID" json:"robot"`
	Product Products `gorm:"foreignKey:ProductID;references:ID" json:"product"`
}

// AiPrediction maps a stored inventory forecast for a product.
type AiPrediction struct {
	ID                uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID         string    `gorm:"type:varchar(50);not null" json:"-"`
	PredictionDate    time.Time `gorm:"type:date;not null" json:"prediction_date"`
	DaysUntilStockout int       `gorm:"not null" json:"days_until_stockout"`
	RecommendedOrder  int       `gorm:"not null" json:"recommended_order"`
	ConfidenceScore   float64   `gorm:"type:decimal(3,2)" json:"confidence_score"`
	CreatedAt         time.Time `gorm:"type:timestamptz;default:now()" json:"created_at"`

	AIPredictionProduct Products `gorm:"foreignKey:ProductID;references:ID;" json:"product"`
}

// TableName returns the database table name used by GORM.
func (InventoryHistory) TableName() string {
	return "inventory_history"
}

// TableName returns the database table name used by GORM.
func (Users) TableName() string {
	return "users"
}

// TableName returns the database table name used by GORM.
func (Products) TableName() string {
	return "products"
}

// TableName returns the database table name used by GORM.
func (Robots) TableName() string {
	return "robots"
}

// TableName returns the database table name used by GORM.
func (AiPrediction) TableName() string {
	return "ai_predictions"
}
