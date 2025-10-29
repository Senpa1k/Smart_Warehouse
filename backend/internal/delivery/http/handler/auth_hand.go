package handler

import (
	"net/http"

	"github.com/Senpa1k/Smart_Warehouse/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) signUp(c *gin.Context) { // not task
	var input models.Users

	if err := c.BindJSON(&input); err != nil {
		NewResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		NewResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Print("sign-up successfuly")
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) login(c *gin.Context) {
	var input models.Users

	if err := c.Bind(&input); err != nil {
		NewResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	token, in, err := h.services.Authorization.GetUser(input.Email, input.PasswordHash)
	if err != nil {
		NewResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if in.ID == 0 {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error":   "invalid date",
			"message": "неверный email или пароль",
		})
		return
	}

	logrus.Print("sign-up successfuly")
	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":   in.ID,
			"name": in.Name,
			"role": in.Role,
		},
	})
}
