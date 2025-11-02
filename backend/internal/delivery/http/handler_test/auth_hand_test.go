package test_handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Senpa1k/Smart_Warehouse/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignUp(t *testing.T) {
	mocks := NewMockServices()
	h := createTestHandler(mocks)

	router := setupTestRouter()
	router.POST("/signup", h.SignUp)

	tests := []struct {
		name           string
		requestBody    models.Users
		mockSetup      func()
		expectedStatus int
	}{
		{
			name: "successful registration",
			requestBody: models.Users{
				Email:        "test@example.com",
				PasswordHash: "password123",
				Name:         "Test User",
			},
			mockSetup: func() {
				mocks.Authorization.On("CreateUser", mock.AnythingOfType("models.Users")).Return(1, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "short password",
			requestBody: models.Users{
				Email:        "test@example.com",
				PasswordHash: "short",
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service error",
			requestBody: models.Users{
				Email:        "test@example.com",
				PasswordHash: "password123",
			},
			mockSetup: func() {
				mocks.Authorization.On("CreateUser", mock.AnythingOfType("models.Users")).Return(0, errors.New("user already exists"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestLogin(t *testing.T) {
	mocks := NewMockServices()
	h := createTestHandler(mocks)

	router := setupTestRouter()
	router.POST("/login", h.Login)

	t.Run("successful login", func(t *testing.T) {
		userInput := models.Users{
			Email:        "test@example.com",
			PasswordHash: "password123",
		}
		userResponse := models.Users{
			ID:   1,
			Name: "Test User",
			Role: "user",
		}

		mocks.Authorization.On("GetUser", userInput.Email, userInput.PasswordHash).Return("jwt-token", userResponse, nil)

		body, _ := json.Marshal(userInput)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "jwt-token", resp["token"])

		userData := resp["user"].(map[string]interface{})
		assert.Equal(t, float64(1), userData["id"])
		assert.Equal(t, "Test User", userData["name"])
		assert.Equal(t, "user", userData["role"])
	})

	t.Run("invalid credentials", func(t *testing.T) {
		userInput := models.Users{
			Email:        "wrong@example.com",
			PasswordHash: "wrongpassword",
		}
		emptyUser := models.Users{ID: 0}

		mocks.Authorization.On("GetUser", userInput.Email, userInput.PasswordHash).Return("", emptyUser, nil)

		body, _ := json.Marshal(userInput)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		userInput := models.Users{
			Email:        "test@example.com",
			PasswordHash: "password123",
		}

		mocks.Authorization.On("GetUser", userInput.Email, userInput.PasswordHash).Return("", models.Users{}, errors.New("database error"))

		body, _ := json.Marshal(userInput)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
