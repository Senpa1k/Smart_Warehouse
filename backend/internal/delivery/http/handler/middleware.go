package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	robotCtx            = "robotId"
)

// UserIdentity validates a user token and stores the user ID in the request context.
func (h *Handler) UserIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	var token string
	if header == "" {
		token = c.Query("token")
		if token == "" {
			NewResponseError(c, http.StatusUnauthorized, "empty auth header or token")
			return
		}
	} else {
		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			NewResponseError(c, http.StatusUnauthorized, "invalid number of auth")
			return
		}
		token = headerParts[1]
	}

	userID, err := h.services.ParseToken(token)
	if err != nil {
		NewResponseError(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userCtx, userID)
}

// RobotIdentity validates a robot authorization header and stores its ID.
func (h *Handler) RobotIdentity(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		NewResponseError(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		NewResponseError(c, http.StatusUnauthorized, "invalid number of auth")
		return
	}

	_, robotID, ok := strings.Cut(headerParts[1], "_")
	if !ok || robotID == "" {
		NewResponseError(c, http.StatusUnauthorized, "invalid robot authorization")
		return
	}
	if !h.services.CheckID(robotID) {
		NewResponseError(c, http.StatusUnauthorized, fmt.Errorf("robot with id=%s does not exist", robotID).Error())
		return
	}
	c.Set(robotCtx, robotID)
}

// WebsocketIdentity checks the headers required for a WebSocket upgrade.
func (h *Handler) WebsocketIdentity(c *gin.Context) {
	if c.GetHeader("Connection") != "Upgrade" {
		NewResponseError(c, http.StatusBadRequest, fmt.Errorf("there is not header Connections").Error())
	}
	if c.GetHeader("Upgrade") != "websocket" {
		NewResponseError(c, http.StatusBadRequest, fmt.Errorf("there is not header Upgrade").Error())
	}
	if c.Request.Header.Get("Sec-WebSocket-Version") == "" {
		c.Request.Header.Set("Sec-WebSocket-Version", "13")
	}

	if c.Request.Header.Get("Sec-WebSocket-Key") == "" {
		NewResponseError(c, http.StatusBadRequest, fmt.Errorf("there is not header sec--key").Error())
	}
}

// RateLimitMiddleware limits requests by client IP and route through Redis.
func (h *Handler) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Если Redis не доступен - пропускаем
		if h.services.Redis == nil {
			c.Next()
			return
		}

		// Создаем ключ на основе IP и пути
		clientIP := c.ClientIP()
		key := fmt.Sprintf("rate:%s:%s", clientIP, c.Request.URL.Path)

		// Проверяем лимит: 100 запросов в минуту
		allowed, err := h.services.Redis.CheckRateLimit(key, 100, time.Minute)
		if err != nil {
			logrus.Errorf("Rate limit error: %v", err)
			c.Next()
			return
		}

		if !allowed {
			NewResponseError(c, http.StatusTooManyRequests, "Слишком много запросов. Попробуйте позже.")
			c.Abort()
			return
		}

		c.Next()
	}
}
