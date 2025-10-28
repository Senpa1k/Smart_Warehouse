package entities

import "time"

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

type AIRequest struct {
	PeriodDays int `json:"period_days" binding:"required"`
	Categories []string `json:"categories" binding:"required"`
}

type AIResponse struct {
	Predictions []struct{
		ProductID string `json:"product_id"`
		PredictionDate string `json:"prediction_date"`
		DaysUntilStockout int `json:"days_until_stockout"`
		RecommendedOrder int `json:"recommended_order"`
		ConfidenceScore float64 `json:"confidence_score"`
	} `json:"predictions"`
	Confidence float64 `json:"confidience"`
}