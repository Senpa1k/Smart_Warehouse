package handler

import (
	"net/http"

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

	h.services.WebsocketDashBoard.RunStream(conn)
	logrus.Print("вебсокет закрыт")
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

	c.JSON(http.StatusOK, dash)
}
