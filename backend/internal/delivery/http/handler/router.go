package handler

import (
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/service"
	"github.com/gin-contrib/cors"
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

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost", "http://localhost:5173"}, // твой фронтенд (Vite, React и т.д.)
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/sign-up", h.signUp)
			auth.POST("/login", h.login)
		}

		robots := api.Group("/robots", h.robotIdentity)
		{
			robots.POST("/data", h.robots)
		}
		ws := api.Group("/ws", h.userIdentity, h.websocketIdentity)
		{
			ws.GET("/dashboard", h.websocketDashBoard)
		}
		inventory := api.Group("/inventory", h.userIdentity)
		{
			inventory.POST("/import", h.importInventory)
			inventory.GET("/history", h.exportInventoryHistory)
		}

		export := api.Group("/export", h.userIdentity)
		{
			export.GET("/excel", h.exportExcel)
			export.GET("/pdf", h.exportPDF)
		}
		dashboard := api.Group("/dashboard", h.userIdentity)
		{
			dashboard.GET("/current", h.getDashInfo)
		}
		ai := api.Group("/ai", h.userIdentity)
		{
			ai.POST("/predict", h.AIRequest)
		}
	}

	return router
}
