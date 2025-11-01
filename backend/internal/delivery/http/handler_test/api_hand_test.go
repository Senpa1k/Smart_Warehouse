package test_handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/delivery/http/handler"
	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// setupTestRouter создает тестовый роутер
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// createTestHandler создает handler с мок сервисами
func createTestHandler(mocks *MockServices) *handler.Handler {
	services := &service.Service{
		Robot:              mocks.Robot,
		DashBoard:          mocks.DashBoard,
		WebsocketDashBoard: mocks.WebsocketDashBoard,
		AI:                 mocks.AI,
		Inventory:          mocks.Inventory,
		Redis:              mocks.Redis,
		Authorization:      mocks.Authorization,
	}

	return handler.NewHandler(services)
}

func TestRobotsEndpoint(t *testing.T) {
	mocks := NewMockServices()
	h := createTestHandler(mocks)

	router := setupTestRouter()
	router.POST("/robots", func(c *gin.Context) {
		c.Set("robotId", "test-robot")
		h.Robots(c)
	})

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func()
		expectedStatus int
	}{
		{
			name: "successful data submission",
			requestBody: entities.RobotsData{
				RobotId:   "test-robot",
				Timestamp: time.Now(),
				Location: struct {
					Zone  string `json:"zone" binding:"required"`
					Row   int    `json:"row" binding:"required"`
					Shelf int    `json:"shelf" binding:"required"`
				}{
					Zone:  "A",
					Row:   1,
					Shelf: 2,
				},
				ScanResults: []entities.ScanResults{
					{
						ProductId:   "prod1",
						ProductName: "Product 1",
						Quantity:    10,
						Status:      "ok",
					},
				},
				BatteryLevel:   80,
				NextCheckpoint: "checkpoint1",
			},
			mockSetup: func() {
				mocks.Robot.On("AddData", mock.AnythingOfType("entities.RobotsData")).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid data format",
			requestBody: map[string]interface{}{
				"invalid": "data",
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service error",
			requestBody: entities.RobotsData{
				RobotId:   "test-robot",
				Timestamp: time.Now(),
				Location: struct {
					Zone  string `json:"zone" binding:"required"`
					Row   int    `json:"row" binding:"required"`
					Shelf int    `json:"shelf" binding:"required"`
				}{
					Zone:  "A",
					Row:   1,
					Shelf: 2,
				},
				ScanResults:    []entities.ScanResults{},
				BatteryLevel:   80,
				NextCheckpoint: "checkpoint1",
			},
			mockSetup: func() {
				mocks.Robot.On("AddData", mock.AnythingOfType("entities.RobotsData")).Return(errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/robots", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mocks.Robot.AssertExpectations(t)
		})
	}
}

func TestGetDashInfo(t *testing.T) {
	mocks := NewMockServices()
	h := createTestHandler(mocks)

	router := setupTestRouter()
	router.GET("/dashboard", func(c *gin.Context) {
		c.Set("userId", 1)
		h.GetDashInfo(c)
	})

	t.Run("successful dashboard info", func(t *testing.T) {
		mocks.DashBoard.On("GetDashInfo", mock.AnythingOfType("*entities.DashInfo")).Return(nil).Run(func(args mock.Arguments) {
			dash := args.Get(0).(*entities.DashInfo)
			dash.Statistics.TotalRobots = 5
			dash.Statistics.ActiveRobots = 3
			dash.Statistics.AvgBattery = 85
		})

		req, _ := http.NewRequest("GET", "/dashboard", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response entities.DashInfo
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, 5, response.Statistics.TotalRobots)
		assert.Equal(t, 3, response.Statistics.ActiveRobots)
	})

	t.Run("dashboard service error", func(t *testing.T) {
		mocks.DashBoard.On("GetDashInfo", mock.AnythingOfType("*entities.DashInfo")).Return(errors.New("service error"))

		req, _ := http.NewRequest("GET", "/dashboard", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestAIRequest(t *testing.T) {
	mocks := NewMockServices()
	h := createTestHandler(mocks)

	router := setupTestRouter()
	router.POST("/ai-predict", func(c *gin.Context) {
		c.Set("userId", 1)
		h.AIRequest(c)
	})

	t.Run("successful AI prediction", func(t *testing.T) {
		request := entities.AIRequest{
			PeriodDays: 7,
			Categories: []string{"electronics", "books"},
		}
		response := &entities.AIResponse{
			Predictions: []entities.Predictions{
				{
					ProductID:         "prod1",
					PredictionDate:    "01.01.2024",
					DaysUntilStockout: 5,
					RecommendedOrder:  100,
					ConfidenceScore:   0.95,
				},
			},
			Confidence: 0.95,
		}

		mocks.AI.On("Predict", request).Return(response, nil)

		body, _ := json.Marshal(request)
		req, _ := http.NewRequest("POST", "/ai-predict", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NotNil(t, resp["predictions"])
		assert.NotNil(t, resp["confidence"])
	})
}

func TestExportExcel(t *testing.T) {
	mocks := NewMockServices()
	h := createTestHandler(mocks)

	router := setupTestRouter()
	router.GET("/export", func(c *gin.Context) {
		c.Set("userId", 1)
		h.ExportExcel(c)
	})

	t.Run("successful export", func(t *testing.T) {
		excelData := []byte("fake excel data")
		productIDs := []string{"1", "2", "3"}

		mocks.Inventory.On("ExportExcel", productIDs).Return(excelData, nil)

		req, _ := http.NewRequest("GET", "/export?ids=1,2,3", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", w.Header().Get("Content-Type"))
		assert.Equal(t, "attachment; filename=inventory.xlsx", w.Header().Get("Content-Disposition"))
		assert.Equal(t, excelData, w.Body.Bytes())
	})

	t.Run("missing ids parameter", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/export", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetRobotsStatus(t *testing.T) {
	mocks := NewMockServices()
	h := createTestHandler(mocks)

	router := setupTestRouter()
	router.GET("/robots-status", h.GetRobotsStatus)

	t.Run("redis available", func(t *testing.T) {
		// Mock responses for each robot
		for _, robotID := range []string{"RB-001", "RB-002", "RB-003", "RB-004", "RB-005"} {
			mocks.Redis.On("IsRobotOnline", robotID).Return(true, nil)
			mocks.Redis.On("GetRobotBattery", robotID).Return(85, nil)
			mocks.Redis.On("GetRobotStatus", robotID).Return("working", nil)
		}

		req, _ := http.NewRequest("GET", "/robots-status", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(5), resp["online_robots"])
		assert.Equal(t, float64(5), resp["total_robots"])
		assert.Equal(t, float64(85), resp["avg_battery"])
	})

	t.Run("redis unavailable", func(t *testing.T) {
		servicesNoRedis := &service.Service{
			Redis: nil,
			// остальные сервисы могут быть nil или моками
		}
		hNoRedis := handler.NewHandler(servicesNoRedis)

		router.GET("/robots-status-no-redis", hNoRedis.GetRobotsStatus)

		req, _ := http.NewRequest("GET", "/robots-status-no-redis", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "Redis не доступен, статусы в реальном времени недоступны", resp["message"])
	})
}
