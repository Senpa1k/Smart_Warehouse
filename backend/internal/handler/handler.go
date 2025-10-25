package handler

import (
	"github.com/Senpa1k/Smart_Warehouse/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sing-up", h.singUp)
		auth.POST("/sing-ip", h.singIn)
	}

	// api := router.Group("/api")
	// {

	// }

	return router
}
