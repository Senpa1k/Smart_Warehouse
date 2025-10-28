package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HistoryQuery struct {
	From   string `form:"from" binding:"required"`
	To     string `form:"to" binding:"required"`
	Zone   string `form:"zone"`
	Status string `form:"status"`
	Limit  int    `form:"limit"`
	Offset int    `form:"offset"`
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

	var query HistoryQuery
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
