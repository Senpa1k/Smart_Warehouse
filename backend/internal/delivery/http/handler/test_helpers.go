package handler

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/models"
)

func CreateTestRobotData() entities.RobotsData {
	return entities.RobotsData{
		RobotId:   "RB-001",
		Timestamp: time.Now(),
		Location: struct {
			Zone  string `json:"zone" binding:"required"`
			Row   int    `json:"row" binding:"required"`
			Shelf int    `json:"shelf" binding:"required"`
		}{
			Zone:  "A",
			Row:   4,
			Shelf: 2,
		},
		ScanResults: []entities.ScanResults{
			{
				ProductId:   "TEL-4567",
				ProductName: "Роутер RT-AC68U",
				Quantity:    45,
				Status:      "OK",
			},
		},
		BatteryLevel:   85,
		NextCheckpoint: "A-5-1",
	}
}

func CreateTestAIRequest() entities.AIRequest {
	return entities.AIRequest{
		PeriodDays: 7,
		Categories: []string{"network"},
	}
}

func CreateTestUser() models.Users {
	return models.Users{
		Email:        "operator@rtk.ru",
		PasswordHash: "password123",
		Name:         "Ёж оператор",
		Role:         "operator",
	}
}

func CreateJSONBody(data interface{}) *bytes.Buffer {
	body, _ := json.Marshal(data)
	return bytes.NewBuffer(body)
}

func CreateTestHistoryResponse() *entities.HistoryResponse {
	return &entities.HistoryResponse{
		Total: 2,
		Items: []models.InventoryHistory{
			{
				ID:        1,
				ProductID: "TEL-4567",
				Quantity:  45,
				Zone:      "A4",
				Status:    "OK",
			},
			{
				ID:        2,
				ProductID: "TEL-8901",
				Quantity:  12,
				Zone:      "B2",
				Status:    "LOW_STOCK",
			},
		},
		Pagination: struct {
			Limit  int `json:"limit"`
			Offset int `json:"offset"`
		}{
			Limit:  50,
			Offset: 0,
		},
	}
}

func CreateTestAIResponse() *entities.AIResponse {
	return &entities.AIResponse{
		Predictions: []entities.Predictions{
			{
				ProductID:         "TEL-4567",
				PredictionDate:    "2024-01-20",
				DaysUntilStockout: 3,
				RecommendedOrder:  100,
				ConfidenceScore:   0.85,
			},
		},
		Confidence: 0.85,
	}
}

func CreateTestUserResponse() *models.Users {
	return &models.Users{
		ID:           1,
		Email:        "operator@rtk.ru",
		Name:         "Ёж оператор",
		Role:         "operator",
		PasswordHash: "hashed_password",
	}
}
