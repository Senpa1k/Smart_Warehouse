package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUserIdentity_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := NewTestHandler()

	mockService.On("ParseToken", "valid_jwt_token").Return(uint(1), nil).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer valid_jwt_token")

	handler.userIdentity(c)

	userID, exists := c.Get(userCtx)
	assert.True(t, exists)
	assert.Equal(t, uint(1), userID)

	mockService.AssertExpectations(t)
}

func TestUserIdentity_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := NewTestHandler()

	mockService.On("ParseToken", "invalid_token").Return(uint(0), assert.AnError).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer invalid_token")

	handler.userIdentity(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "error")

	mockService.AssertExpectations(t)
}

func TestRobotIdentity_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := NewTestHandler()

	mockService.On("CheckId", "RB-001").Return(true).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/robots/data", nil)
	c.Request.Header.Set("Authorization", "Bearer robot_RB-001")

	handler.robotIdentity(c)

	robotID, exists := c.Get(robotCtx)
	assert.True(t, exists)
	assert.Equal(t, "RB-001", robotID)

	mockService.AssertExpectations(t)
}

func TestRobotIdentity_InvalidRobot(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := NewTestHandler()

	mockService.On("CheckId", "RB-999").Return(false).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/robots/data", nil)
	c.Request.Header.Set("Authorization", "Bearer robot_RB-999")

	handler.robotIdentity(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "robot with id=RB-999 does not exist")

	mockService.AssertExpectations(t)
}
