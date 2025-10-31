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
	CheckOrigin: func(r *http.Request) bool {
		return true // для разрабы
	},
	WriteBufferSize: 1024,
	ReadBufferSize:  10,
}

func (h *Handler) robots(c *gin.Context) {
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

	err := h.services.Robot.AddData(rd)
	if err != nil {
		NewResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Print("data received successfuly")
	c.JSON(http.StatusOK, gin.H{
		"status":     "received",
		"message_id": "tmp_message_id",
	})
}

func (h *Handler) websocketDashBoard(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.Error("upgrade error with socket")
		return
	}
	defer conn.Close()

	// ✅ НОВОЕ: Подписываемся на Redis channel
	if h.services.Redis != nil {
		go h.handleRedisSubscriptions(conn)
	}

	// Старая логика
	h.services.WebsocketDashBoard.RunStream(conn)
	logrus.Print("вебсокет закрыт")
}

// ✅ НОВАЯ ФУНКЦИЯ: Обработка Redis подписок
func (h *Handler) handleRedisSubscriptions(conn *websocket.Conn) {
	ctx := context.Background()

	// Подписываемся на канал robot_updates
	pubsub := h.services.Redis.Subscribe("robot_updates")
	defer pubsub.Close()

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

		logrus.Info("📨 Sent Redis message to WebSocket client")
	}
}

func (h *Handler) getDashInfo(c *gin.Context) {
	_, ok := c.Get(userCtx)
	if !ok {
		NewResponseError(c, http.StatusInternalServerError, "robot id not found")
		return
	}

	var dash entities.DashInfo = entities.DashInfo{}
	if err := h.services.DashBoard.GetDashInfo(&dash); err != nil {
		NewResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Print("dashinfo update successfuly")
	c.JSON(http.StatusOK, dash)
}

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

	res, err := h.services.AI.Predict(air)
	if err != nil {
		NewResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Print("ai request send successfuly")
	c.JSON(http.StatusOK, gin.H{
		"predictions": res.Predictions,
		"confidence":  res.Confidence,
	})
}

func (h *Handler) exportExcel(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok {
		NewResponseError(c, http.StatusUnauthorized, "user not authenticated")
		return
	}
	_ = userID

	productIdStr := c.Query("ids")
	if productIdStr == "" {
		NewResponseError(c, http.StatusBadRequest, "ids query parameter is required")
		return
	}

	productIDs := strings.Split(productIdStr, ",")
	for i, id := range productIDs {
		productIDs[i] = strings.TrimSpace(id)
	}

	exelFile, err := h.services.Inventory.ExportExcel(productIDs)
	if err != nil {
		NewResponseError(c, http.StatusInternalServerError, "export failed: "+err.Error())
		return
	}
	c.Header("Content-Disposition", "attachment; filename=inventory.xlsx")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", exelFile)
}

func (h *Handler) importInventory(c *gin.Context) {
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
	defer file.Close()

	result, err := h.services.Inventory.ImportCSV(file)
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

	historyData, err := h.services.Inventory.GetHistory(query.From, query.To, query.Zone, query.Status, query.Limit, query.Offset)
	if err != nil {
		NewResponseError(c, http.StatusInternalServerError, "failed to get history: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, historyData)
}

// Получение статусов всех роботов
func (h *Handler) getRobotsStatus(c *gin.Context) {
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

	statuses := make(map[string]interface{})
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

		statuses[robotID] = map[string]interface{}{
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
