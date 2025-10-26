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
	if header == "" {
		NewResponseError(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		NewResponseError(c, http.StatusUnauthorized, "invalid number of auth")
		return
	}

	userID, err := h.services.Authorization.ParseToken(headerParts[1])
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
