package test_handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserIdentityMiddleware(t *testing.T) {
	mocks := NewMockServices()
	h := createTestHandler(mocks)

	router := setupTestRouter()
	router.Use(h.UserIdentity)
	router.GET("/protected", func(c *gin.Context) {
		userID, exists := c.Get("userId")
		if exists {
			c.JSON(http.StatusOK, gin.H{"user_id": userID})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		}
	})

	tests := []struct {
		name           string
		authHeader     string
		mockSetup      func()
		expectedStatus int
	}{
		{
			name:       "valid token in header",
			authHeader: "Bearer valid-token",
			mockSetup: func() {
				mocks.Authorization.On("ParseToken", "valid-token").Return(1, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:       "valid token in query",
			authHeader: "",
			mockSetup: func() {
				mocks.Authorization.On("ParseToken", "query-token").Return(1, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "no auth provided",
			authHeader:     "",
			mockSetup:      func() {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid token format",
			authHeader:     "InvalidFormat",
			mockSetup:      func() {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:       "invalid token",
			authHeader: "Bearer invalid-token",
			mockSetup: func() {
				mocks.Authorization.On("ParseToken", "invalid-token").Return(0, errors.New("invalid token"))
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req, _ := http.NewRequest("GET", "/protected", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			} else if tt.name == "valid token in query" {
				req.URL.RawQuery = "token=query-token"
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestRobotIdentityMiddleware(t *testing.T) {
	mocks := NewMockServices()
	h := createTestHandler(mocks)

	router := setupTestRouter()
	router.Use(h.RobotIdentity)
	router.POST("/robot-data", func(c *gin.Context) {
		robotID, exists := c.Get("robotId")
		if exists {
			c.JSON(http.StatusOK, gin.H{"robot_id": robotID})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "robot not found"})
		}
	})

	t.Run("valid robot authorization", func(t *testing.T) {
		mocks.Robot.On("CheckId", "test-robot").Return(true)

		req, _ := http.NewRequest("POST", "/robot-data", nil)
		req.Header.Set("Authorization", "Bearer Robot_test-robot")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid robot id", func(t *testing.T) {
		mocks.Robot.On("CheckId", "invalid-robot").Return(false)

		req, _ := http.NewRequest("POST", "/robot-data", nil)
		req.Header.Set("Authorization", "Robot_invalid-robot")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestRateLimitMiddleware(t *testing.T) {
	mocks := NewMockServices()
	h := createTestHandler(mocks)

	router := setupTestRouter()
	router.Use(h.RateLimitMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	t.Run("within rate limit", func(t *testing.T) {
		mocks.Redis.On("CheckRateLimit", mock.Anything, 100, mock.Anything).Return(true, nil)

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("exceeds rate limit", func(t *testing.T) {
		mocks.Redis.On("CheckRateLimit", mock.Anything, 100, mock.Anything).Return(false, nil)

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTooManyRequests, w.Code)
	})

	t.Run("redis error", func(t *testing.T) {
		mocks.Redis.On("CheckRateLimit", mock.Anything, 100, mock.Anything).Return(false, errors.New("redis error"))

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code) // Should proceed on error
	})
}
