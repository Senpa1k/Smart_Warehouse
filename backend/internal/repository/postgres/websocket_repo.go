package postgres

import (
	"fmt"
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

	enti.Type = "inventory_alert"
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

func (w *WebsocketDashBoardPostgres) InventoryAlertPredict(enti *entities.InventoryAlert, predict entities.Predictions) error {
	var product models.Products
	if err := w.db.Where("id = ?", predict.ProductID).First(&product).Error; err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	// Parse predicted stockout date to calculate days until stockout
	stockoutDate, err := time.Parse("2006-01-02", predict.PredictedStockoutDate)
	if err != nil {
		return fmt.Errorf("failed to parse stockout date: %w", err)
	}

	daysUntilStockout := int(time.Until(stockoutDate).Hours() / 24)
	if daysUntilStockout < 0 {
		daysUntilStockout = 0
	}

	alertType := "predicted"
	message := ""

	switch {
	case daysUntilStockout <= 2:
		message = fmt.Sprintf("КРИТИЧЕСКИЙ УРОВЕНЬ! Товар закончится через %d дней", daysUntilStockout)
	case daysUntilStockout <= 7:
		message = fmt.Sprintf("Рекомендуется пополнение. Товар закончится через %d дней", daysUntilStockout)
	default:
		return fmt.Errorf("everything norm")
	}

	status := "OK"
	if predict.CurrentStock <= product.MinStock {
		status = "CRITICAL"
	} else if predict.CurrentStock <= product.OptimalStock/2 {
		status = "LOW_STOCK"
	}

	enti.Type = "inventory_alert"
	enti.Data.ProductId = predict.ProductID
	enti.Data.ProductName = predict.ProductName
	enti.Data.CurrentQuantity = predict.CurrentStock
	enti.Data.Zone = ""
	enti.Data.Row = 0
	enti.Data.Shelf = 0
	enti.Data.Status = status
	enti.Data.AlterType = alertType
	enti.Data.Timestamp = time.Now()
	enti.Data.Message = fmt.Sprintf(
		"%s. Рекомендуемый заказ: %d единиц. Уверенность прогноза: %.1f%%",
		message,
		predict.RecommendedOrderQuantity,
		predict.ConfidenceScore*100,
	)

	return nil
}
