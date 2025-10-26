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

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			// auth.POST("/sing-up", h.singUp)
			auth.POST("/login", h.login)
		}

		robots := api.Group("/robots", h.robotIdentity)
		{
			robots.POST("/data", h.robots)
		}

		dashboard := api.Group("/dashboard", h.userIdentity)
		{
			dashboard.GET("/current", h.currentDashBoard)
		}

	}

	return router
}
