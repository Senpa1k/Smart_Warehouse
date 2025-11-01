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

	router.Use(h.RateLimitMiddleware())

	// разрешаем cors
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost", "http://localhost:5173"},
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
			auth.POST("/sign-up", h.SignUp)
			auth.POST("/login", h.Login)
		}

		robots := api.Group("/robots", h.RobotIdentity)
		{
			robots.POST("/data", h.Robots)
		}
		ws := api.Group("/ws", h.UserIdentity, h.WebsocketIdentity)
		{
			ws.GET("/dashboard", h.WebsocketDashBoard)
		}
		inventory := api.Group("/inventory", h.UserIdentity)
		{
			inventory.POST("/import", h.ImportInventory)
			inventory.GET("/history", h.exportInventoryHistory)
		}

		export := api.Group("/export", h.UserIdentity)
		{
			export.GET("/excel", h.ExportExcel)
		}
		dashboard := api.Group("/dashboard", h.UserIdentity)
		{
			dashboard.GET("/current", h.GetDashInfo)
		}
		ai := api.Group("/ai", h.UserIdentity)
		{
			ai.POST("/predict", h.AIRequest)
		}

		monitoring := api.Group("/monitoring", h.UserIdentity)
		{
			monitoring.GET("/robots/status", h.GetRobotsStatus)
		}
	}

	return router
}
