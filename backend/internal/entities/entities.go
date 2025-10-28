package entities

import (
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/models"
)

type RobotsData struct {
	RobotId   string    `json:"robot_id" binding:"required"`
	Timestamp time.Time `json:"timestamp" binding:"required"`
	Location  struct {
		Zone  string `json:"zone" binding:"required"`
		Row   int    `json:"row" binding:"required"`
		Shelf int    `json:"shelf" binding:"required"`
	} `json:"location" binding:"required"`
	ScanResults    []ScanResults `json:"scan_results" binding:"required"`
	BatteryLevel   int           `json:"battery_level" binding:"required"`
	NextCheckpoint string        `json:"next_checkpoint" binding:"required"`
}

type ScanResults struct {
	ProductId   string `json:"product_id" binding:"required"`
	ProductName string `json:"product_name" binding:"required"`
	Quantity    int    `json:"quantity" binding:"required"`
	Status      string `json:"status" binding:"required"`
}

type DashInfo struct {
	ListRobots []models.Robots              `json:"robots"`
	ListScans  [100]models.InventoryHistory `json:"recent_scans"`
	Statistics Statistics                   `json:"statistics"`
}

type Statistics struct {
	TotalCheck        int `json:"total_check"`
	UniqueProducts    int `json:"unique_products"`
	FindDiscrepancies int `json:"find _discrepancies"`
	AverageTime       int `json:"average_time"`
}

type InventoryAlert struct {
	Type string `json:"type"`
	Data struct {
		ProductId       string    `json:"product_id"`
		ProductName     string    `json:"product_name"`
		CurrentQuantity int       `json:"current_quantity"`
		Zone            string    `json:"zone"`
		Row             int       `json:"row"`
		Shelf           int       `json:"shelf"`
		Status          string    `json:"status"`
		AlterType       string    `json:"alter_type"`
		Timestamp       time.Time `json:"timestamp"`
		Message         string    `json:"message"`
	} `json:"data"`
}

type UpdateRobot struct {
	ID           string    `gorm:"primaryKey;type:varchar(50)" json:"id"`
	Status       string    `gorm:"size:50;default:active" json:"status"`
	BatteryLevel int       `gorm:"type:int" json:"battery_level"`
	LastUpdate   time.Time `gorm:"type:timestamptz" json:"last_update"`
	CurrentZone  string    `gorm:"size:10" json:"current_zone"`
	CurrentRow   int       `gorm:"type:int" json:"current_row"`
	CurrentShelf int       `gorm:"type:int" json:"current_shelf"`
}

type AIRequest struct {
	PeriodDays int      `json:"period_days" binding:"required"`
	Categories []string `json:"categories" binding:"required"`
}

type AIResponse struct {
	Predictions []struct {
		ProductID         string  `json:"product_id"`
		PredictionDate    string  `json:"prediction_date"`
		DaysUntilStockout int     `json:"days_until_stockout"`
		RecommendedOrder  int     `json:"recommended_order"`
		ConfidenceScore   float64 `json:"confidence_score"`
	} `json:"predictions"`
	Confidence float64 `json:"confidience"`
}
