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
	ListRobots []*models.Robots              `json:"robots"`
	ListScans  [100]*models.InventoryHistory `json:"recent_scans"`
	Statistics struct {
		TotalCheck        int `json:"total_check"`
		UniqueProducts    int `json:"unique_products"`
		FindDiscrepancies int `json:"find _discrepancies"`
		AverageTime       int `json:"average_time"`
	} `json:"statistics"`
}
