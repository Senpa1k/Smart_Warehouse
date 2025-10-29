package postgres

import (
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/models"
	"gorm.io/gorm"
)

type WebsocketDashBoardPostgres struct {
	db *gorm.DB
}

func NewWebsocketDashBoardPostgres(db *gorm.DB) *WebsocketDashBoardPostgres {
	return &WebsocketDashBoardPostgres{db: db}
}

func (w *WebsocketDashBoardPostgres) InventoryAlertScanned(enti *entities.InventoryAlert, timestemp time.Time, idProduct string) error {
	var inventoryHistory models.InventoryHistory

	err := w.db.Where("scanned_at = ? and product_id = ?", timestemp, idProduct).First(&inventoryHistory).Error
	if err != nil {
		return err
	}

	var product string
	err = w.db.Model(models.Products{}).Where("id = ?", inventoryHistory.ProductID).Select("name").Scan(&product).Error
	if err != nil {
		return err
	}

	enti.Data.ProductId = inventoryHistory.ProductID
	enti.Data.ProductName = product
	enti.Data.CurrentQuantity = inventoryHistory.Quantity
	enti.Data.Zone = inventoryHistory.Zone
	enti.Data.Row = inventoryHistory.RowNumber
	enti.Data.Shelf = inventoryHistory.ShelfNumber
	enti.Data.Status = inventoryHistory.Status
	enti.Data.Timestamp = inventoryHistory.ScannedAt
	enti.Data.AlterType = "scanned"
	enti.Data.Message = inventoryHistory.Status + " остаток! Требуется пополнение."

	return nil
}
