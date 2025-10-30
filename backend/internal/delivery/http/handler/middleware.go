package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	robotCtx            = "robotId"
)

func (h *Handler) userIdentity(c *gin.Context) {

	header := c.GetHeader(authorizationHeader)
	token := ""
	if header == "" {
		// Для WebSocket, проверяем query параметр token
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

	userID, err := h.services.Authorization.ParseToken(token)
	if err != nil {
		NewResponseError(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userCtx, userID)
}

func (h *Handler) robotIdentity(c *gin.Context) {
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

	robotID := strings.Split(headerParts[1], "_")[1]
	if !h.services.Robot.CheckId(robotID) {
		NewResponseError(c, http.StatusUnauthorized, fmt.Errorf("robot with id=%s does not exist", robotID).Error())
		return
	}
	c.Set(robotCtx, robotID)
}

func (h *Handler) websocketIdentity(c *gin.Context) {
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
