// Package entities defines warehouse API requests, responses, and events.
package entities

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/models"
)

// ParseCheckpoint decodes a zone-row-shelf checkpoint and rejects malformed coordinates.
func ParseCheckpoint(checkpoint string) (string, int, int, error) {
	parts := strings.Split(checkpoint, "-")
	if len(parts) != 3 || parts[0] == "" {
		return "", 0, 0, fmt.Errorf("invalid checkpoint %q: expected zone-row-shelf", checkpoint)
	}
	row, err := strconv.Atoi(parts[1])
	if err != nil || row < 0 {
		return "", 0, 0, fmt.Errorf("invalid checkpoint row in %q", checkpoint)
	}
	shelf, err := strconv.Atoi(parts[2])
	if err != nil || shelf < 0 {
		return "", 0, 0, fmt.Errorf("invalid checkpoint shelf in %q", checkpoint)
	}
	return parts[0], row, shelf, nil
}

// RobotsData carries a robot position and its inventory scans.
type RobotsData struct {
	RobotID   string    `json:"robot_id" binding:"required"`
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

// ScanResults describes one product observation made by a robot.
type ScanResults struct {
	ProductID   string `json:"product_id" binding:"required"`
	ProductName string `json:"product_name" binding:"required"`
	Quantity    int    `json:"quantity" binding:"required"`
	Status      string `json:"status" binding:"required"`
}

// DashInfo contains the current robots, recent scans, and summary statistics.
type DashInfo struct {
	ListRobots []models.Robots              `json:"robots"`
	ListScans  [100]models.InventoryHistory `json:"recent_scans"`
	Statistics Statistics                   `json:"statistics"`
}

// Statistics summarizes robot activity and inventory alerts.
type Statistics struct {
	ActiveRobots      int `json:"active_robots"`
	TotalRobots       int `json:"total_robots"`
	ItemsCheckedToday int `json:"items_checked_today"`
	CriticalItems     int `json:"critical_items"`
	AvgBattery        int `json:"avg_battery"`
}

// InventoryAlert is a WebSocket event for a low-stock or predicted shortage.
type InventoryAlert struct {
	Type string `json:"type"`
	Data struct {
		ProductID       string    `json:"product_id"`
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

// UpdateRobot is the robot state sent to dashboard WebSocket clients.
type UpdateRobot struct {
	ID           string    `gorm:"primaryKey;type:varchar(50)" json:"id"`
	Status       string    `gorm:"size:50;default:active" json:"status"`
	BatteryLevel int       `gorm:"type:int" json:"battery_level"`
	LastUpdate   time.Time `gorm:"type:timestamptz" json:"last_update"`
	CurrentZone  string    `gorm:"size:10" json:"current_zone"`
	CurrentRow   int       `gorm:"type:int" json:"current_row"`
	CurrentShelf int       `gorm:"type:int" json:"current_shelf"`
}

// AIRequest selects the forecast period and product categories.
type AIRequest struct {
	PeriodDays int      `json:"period_days" binding:"required"`
	Categories []string `json:"categories" binding:"required"`
}

// AIResponse contains inventory forecasts and their overall confidence.
type AIResponse struct {
	Predictions []Predictions `json:"predictions"`
	Confidence  float64       `json:"confidience"`
}

// Predictions describes the forecast and order recommendation for a product.
type Predictions struct {
	ProductID         string  `json:"product_id"`
	ProductName       string  `json:"product_name"`
	PredictionDate    string  `json:"prediction_date"`
	DaysUntilStockout int     `json:"days_until_stockout"`
	RecommendedOrder  int     `json:"recommended_order"`
	ConfidenceScore   float64 `json:"confidence_score"`
}

// ImportResult counts successful and failed CSV records.
type ImportResult struct {
	SuccessCount int      `json:"success_count"`
	FailedCount  int      `json:"failed_count"`
	Errors       []string `json:"errors"`
}

// HistoryResponse contains paginated inventory scan history.
type HistoryResponse struct {
	Total      int64                     `json:"total"`
	Items      []models.InventoryHistory `json:"items"`
	Pagination struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	} `json:"pagination"`
}

// HistoryQuery holds filters and pagination for inventory history.
type HistoryQuery struct {
	From   string `form:"from"`
	To     string `form:"to"`
	Zone   string `form:"zone"`
	Status string `form:"status"`
	Limit  int    `form:"limit"`
	Offset int    `form:"offset"`
}
