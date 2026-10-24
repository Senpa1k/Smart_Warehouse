// Package handler serves warehouse HTTP and WebSocket requests.
package handler

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool {
		return true // для разрабы
	},
	WriteBufferSize: 1024,
	ReadBufferSize:  10,
}

// Robots accepts inventory scans sent by an authenticated robot.
func (h *Handler) Robots(c *gin.Context) {
	_, ok := c.Get(robotCtx)
	if !ok {
		NewResponseError(c, http.StatusInternalServerError, "robot id not found")
		return
	}

	var rd entities.RobotsData
	if err := c.Bind(&rd); err != nil {
		NewResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.AddData(rd)
	if err != nil {
		NewResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Print("data received successfully")
	c.JSON(http.StatusOK, gin.H{
		"status":     "received",
		"message_id": "tmp_message_id",
	})
}

// WebsocketDashBoard upgrades a dashboard request to a WebSocket connection.
func (h *Handler) WebsocketDashBoard(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.Error("upgrade error with socket")
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			logrus.WithError(err).Warn("close WebSocket connection")
		}
	}()

	// Подписываемся на Redis channel
	if h.services.Redis != nil {
		go h.HandleRedisSubscriptions(conn)
	}

	// Старая логика
	h.services.RunStream(conn)
	logrus.Print("вебсокет закрыт")
}

// HandleRedisSubscriptions forwards Redis robot updates to a WebSocket client.
func (h *Handler) HandleRedisSubscriptions(conn *websocket.Conn) {
	ctx := context.Background()

	// Подписываемся на канал robot_updates
	pubsub := h.services.Redis.Subscribe("robot_updates")
	defer func() {
		if err := pubsub.Close(); err != nil {
			logrus.WithError(err).Warn("close Redis subscription")
		}
	}()

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			logrus.Errorf("Redis subscription error: %v", err)
			return
		}

		// Отправляем сообщение через WebSocket
		err = conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
		if err != nil {
			logrus.Errorf("WebSocket send error: %v", err)
			return
		}
		logrus.Info("Sent Redis message to WebSocket client")
	}
}

// GetDashInfo returns the current warehouse dashboard snapshot.
func (h *Handler) GetDashInfo(c *gin.Context) {
	_, ok := c.Get(userCtx)
	if !ok {
		NewResponseError(c, http.StatusInternalServerError, "robot id not found")
		return
	}

	dash := entities.DashInfo{}
	if err := h.services.GetDashInfo(&dash); err != nil {
		NewResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Print("dashinfo update successfully")
	c.JSON(http.StatusOK, dash)
}

// AIRequest requests inventory predictions and writes them to the response.
func (h *Handler) AIRequest(c *gin.Context) {
	_, ok := c.Get(userCtx)
	if !ok {
		NewResponseError(c, http.StatusInternalServerError, "user not found")
		return
	}

	var air entities.AIRequest
	if err := c.Bind(&air); err != nil {
		NewResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.services.Predict(air)
	if err != nil {
		NewResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Print("ai request send successfully")
	c.JSON(http.StatusOK, gin.H{
		"predictions": res.Predictions,
		"confidence":  res.Confidence,
	})
}

// ExportExcel downloads selected inventory scans as an Excel workbook.
func (h *Handler) ExportExcel(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok {
		NewResponseError(c, http.StatusUnauthorized, "user not authenticated")
		return
	}
	_ = userID

	productIDStr := c.Query("ids")
	if productIDStr == "" {
		NewResponseError(c, http.StatusBadRequest, "ids query parameter is required")
		return
	}

	productIDs := strings.Split(productIDStr, ",")
	for i, id := range productIDs {
		productIDs[i] = strings.TrimSpace(id)
	}

	exelFile, err := h.services.ExportExcel(productIDs)
	if err != nil {
		NewResponseError(c, http.StatusInternalServerError, "export failed: "+err.Error())
		return
	}
	c.Header("Content-Disposition", "attachment; filename=inventory.xlsx")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", exelFile)
}

// ImportInventory imports uploaded CSV inventory records.
func (h *Handler) ImportInventory(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok {
		NewResponseError(c, http.StatusUnauthorized, "user not authenticated")
		return
	}
	_ = userID

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		NewResponseError(c, http.StatusBadRequest, "failed to get file: "+err.Error())
		return
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			logrus.WithError(closeErr).Warn("close imported file")
		}
	}()

	result, err := h.services.ImportCSV(file)
	if err != nil {
		NewResponseError(c, http.StatusInternalServerError, "import failed: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) exportInventoryHistory(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok {
		NewResponseError(c, http.StatusUnauthorized, "user not authenticated")
		return
	}
	_ = userID

	var query entities.HistoryQuery
	if err := c.BindQuery(&query); err != nil {
		NewResponseError(c, http.StatusBadRequest, "invalid query parameters: "+err.Error())
		return
	}
	if query.Limit == 0 {
		query.Limit = 50
	}
	if query.Limit > 1000 {
		query.Limit = 1000
	}

	historyData, err := h.services.GetHistory(query.From, query.To, query.Zone, query.Status, query.Limit, query.Offset)
	if err != nil {
		NewResponseError(c, http.StatusInternalServerError, "failed to get history: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, historyData)
}

// GetRobotsStatus reports the online state of the configured demo robots.
func (h *Handler) GetRobotsStatus(c *gin.Context) {
	if h.services.Redis == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "Redis не доступен, статусы в реальном времени недоступны",
			"robots":  []string{},
		})
		return
	}

	// В реальном приложении здесь бы брали список роботов из БД
	// Для демо используем тестовые ID
	robotIDs := []string{"RB-001", "RB-002", "RB-003", "RB-004", "RB-005"}

	statuses := make(map[string]any)
	onlineCount := 0
	totalBattery := 0

	for _, robotID := range robotIDs {
		online, _ := h.services.Redis.IsRobotOnline(robotID)
		battery, _ := h.services.Redis.GetRobotBattery(robotID)
		status, _ := h.services.Redis.GetRobotStatus(robotID)

		if online {
			onlineCount++
			totalBattery += battery
		}

		statuses[robotID] = map[string]any{
			"online":        online,
			"battery_level": battery,
			"status":        status,
			"last_update":   "в реальном времени", // В продакшене хранили бы время
		}
	}

	// Вычисляем среднюю батарею
	avgBattery := 0
	if onlineCount > 0 {
		avgBattery = totalBattery / onlineCount
	}

	c.JSON(http.StatusOK, gin.H{
		"online_robots": onlineCount,
		"total_robots":  len(robotIDs),
		"avg_battery":   avgBattery,
		"robots":        statuses,
		"last_updated":  time.Now().Format("15:04:05"),
	})
}
