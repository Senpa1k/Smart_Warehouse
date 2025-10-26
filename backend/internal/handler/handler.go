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

	router.Group("/api")
	{
		auth := router.Group("/auth")
		{
			// auth.POST("/sing-up", h.singUp)
			auth.POST("/login", h.login)
		}

		// robots := router.Group("/robots")
		// {
		// 	robots.POST("/data", h.robots)
		// }

	}

	return router
}
