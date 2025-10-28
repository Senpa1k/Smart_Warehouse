package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

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
