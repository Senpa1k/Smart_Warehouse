package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		NewResponseError(c, http.StatusUnauthorized, "empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		NewResponseError(c, http.StatusUnauthorized, "invalid number of auth")
	}

	userID, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		NewResponseError(c, http.StatusUnauthorized, err.Error())
	}

	c.Set(userCtx, userID)
}
