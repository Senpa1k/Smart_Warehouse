package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Senpa1k/Smart_Warehouse/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLogin_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := NewTestHandler()

	user := CreateTestUser()
	expectedUser := &models.Users{
		ID:           1,
		Email:        "operator@rtk.ru",
		Name:         "Иван Операторов",
		Role:         "operator",
		PasswordHash: "hashed_password",
	}

	mockService.On("GetUser", "operator@rtk.ru", "password123").Return("jwt_token_here", expectedUser, nil).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("POST", "/api/auth/login", CreateJSONBody(user))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.login(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "jwt_token_here", response["token"])

	userData := response["user"].(map[string]interface{})
	assert.Equal(t, float64(1), userData["id"])
	assert.Equal(t, "Иван Операторов", userData["name"])
	assert.Equal(t, "operator", userData["role"])

	mockService.AssertExpectations(t)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := NewTestHandler()

	user := CreateTestUser()
	emptyUser := &models.Users{ID: 0}

	mockService.On("GetUser", "operator@rtk.ru", "password123").Return("", emptyUser, nil).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("POST", "/api/auth/login", CreateJSONBody(user))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.login(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid date", response["error"])
	assert.Equal(t, "неверный email или пароль", response["message"])

	mockService.AssertExpectations(t)
}

func TestSignUp_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := NewTestHandler()

	user := CreateTestUser()

	mockService.On("CreateUser", user).Return(uint(2), nil).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("POST", "/api/auth/signup", CreateJSONBody(user))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.signUp(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(2), response["id"])

	mockService.AssertExpectations(t)
}
