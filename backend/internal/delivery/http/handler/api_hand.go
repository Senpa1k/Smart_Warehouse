package handler

import (
	"net/http"

	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/gin-gonic/gin"
)

func (h *Handler) robots(c *gin.Context) {
	_, ok := c.Get(userCtx)
	if !ok {
		NewResponseError(c, http.StatusInternalServerError, "robot id not found")
		return
	}

	var rd entities.RobotsData
	if err := c.Bind(&rd); err != nil {
		NewResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.Robot.AddData(rd)
	if err != nil {
		NewResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "received",
		"message_id": "tmp_message_id",
	})
}
