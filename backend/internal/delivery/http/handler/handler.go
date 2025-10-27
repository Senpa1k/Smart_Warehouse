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

		// robots := router.Group("/robots", h.userIdentity)
		robots := api.Group("/robots")
		{
			robots.POST("/data", h.robots)
		}

		inventory := api.Group("/inventory", h.userIdentity)
		{
			inventory.POST("/import", h.importInventory)
			inventory.GET("/history", h.exportInventoryHistory)
		}

		export := api.Group("/export", h.userIdentity)
		{
			export.GET("/excel", h.exportExcel)
		}

	}

	return router
}
