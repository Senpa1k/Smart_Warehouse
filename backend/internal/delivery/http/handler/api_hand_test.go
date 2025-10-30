package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRobots_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := NewTestHandler()

	robotData := CreateTestRobotData()

	mockService.On("AddData", mock.AnythingOfType("entities.RobotsData")).Return(nil).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(robotCtx, "RB-001")

	c.Request = httptest.NewRequest("POST", "/api/robots/data", CreateJSONBody(robotData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.robots(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestRobots_NoRobotContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := NewTestHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.robots(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "robot id not found")
	mockService.AssertNotCalled(t, "AddData")
}

func TestRobots_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := NewTestHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(robotCtx, "RB-001")

	c.Request = httptest.NewRequest("POST", "/api/robots/data", bytes.NewBufferString("invalid json"))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.robots(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "AddData")
}

func TestGetDashInfo_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := NewTestHandler()

	mockService.On("GetDashInfo", mock.AnythingOfType("*entities.DashInfo")).Return(nil).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(userCtx, uint(1))

	handler.getDashInfo(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetDashInfo_NoUserContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := NewTestHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.getDashInfo(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "robot id not found")
	mockService.AssertNotCalled(t, "GetDashInfo")
}

func TestAIRequest_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := NewTestHandler()

	aiRequest := CreateTestAIRequest()
	aiResponse := CreateTestAIResponse()

	mockService.On("Predict", aiRequest).Return(aiResponse, nil).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(userCtx, uint(1))

	c.Request = httptest.NewRequest("POST", "/api/ai/predict", CreateJSONBody(aiRequest))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.AIRequest(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "predictions")
	assert.Contains(t, response, "confidence")

	mockService.AssertExpectations(t)
}

func TestAIRequest_NoUserContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := NewTestHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.AIRequest(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "user not found")
	mockService.AssertNotCalled(t, "Predict")
}

func TestExportInventoryHistory_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := NewTestHandler()

	historyResponse := CreateTestHistoryResponse()

	mockService.On("GetHistory", "2024-01-01", "2024-01-31", "A", "OK", 50, 0).Return(historyResponse, nil).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(userCtx, uint(1))

	c.Request = httptest.NewRequest("GET", "/api/inventory/history?from=2024-01-01&to=2024-01-31&zone=A&status=OK", nil)

	handler.exportInventoryHistory(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response entities.HistoryResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), response.Total)
	assert.Len(t, response.Items, 2)

	mockService.AssertExpectations(t)
}
