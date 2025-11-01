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

// создание оповещений на основе данных инвентаризации
func (w *WebsocketDashBoardPostgres) InventoryAlertScanned(enti *entities.InventoryAlert, timestemp time.Time, idProduct string) error {
	var inventoryHistory models.InventoryHistory

	// получение истории сканирования
	err := w.db.Where("scanned_at = ? and product_id = ?", timestemp, idProduct).First(&inventoryHistory).Error
	if err != nil {
		return err
	}

	// получение названия продукта
	var product string
	err = w.db.Model(models.Products{}).Where("id = ?", inventoryHistory.ProductID).Select("name").Scan(&product).Error
	if err != nil {
		return err
	}

	// формирование оповещения
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


// создание оповещений на основе предективной аналитики
func (w *WebsocketDashBoardPostgres) InventoryAlertPredict(enti *entities.InventoryAlert, predict entities.Predictions) error {
	// получение информации о продукте
	var product models.Products
	if err := w.db.Where("id = ?", predict.ProductID).First(&product).Error; err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	// получение информации о количестве продукта
	var currentQuantity int
	if err := w.db.Table("inventory_history").Where("product_id = ?", predict.ProductID).Select("COALESCE(quantity, 0)").Scan(&currentQuantity).Error; err != nil {
		return fmt.Errorf("failed to get quantity: %w", err)
	}

	alertType := "predicted"
	message := ""

	// определение статуса прогноза
	switch {
	case predict.DaysUntilStockout <= 2:
		message = fmt.Sprintf("КРИТИЧЕСКИЙ УРОВЕНЬ! Товар закончится через %d дней", predict.DaysUntilStockout)
	case predict.DaysUntilStockout <= 7:
		message = fmt.Sprintf("Рекомендуется пополнение. Товар закончится через %d дней", predict.DaysUntilStockout)
	default:
		return fmt.Errorf("everything norm")
	}

	status := "OK"
	if currentQuantity <= product.MinStock {
		status = "CRITICAL"
	} else if currentQuantity <= product.OptimalStock/2 {
		status = "LOW_STOCK"
	}

	// формирование оповещения
	enti.Type = "inventory_alert"
	enti.Data.ProductId = predict.ProductID
	enti.Data.ProductName = product.Name
	enti.Data.CurrentQuantity = currentQuantity
	enti.Data.Zone = ""
	enti.Data.Row = 0
	enti.Data.Shelf = 0
	enti.Data.Status = status
	enti.Data.AlterType = alertType
	enti.Data.Timestamp = time.Now()
	enti.Data.Message = fmt.Sprintf(
		"%s. Рекомендуемый заказ: %d единиц. Уверенность прогноза: %.1f%%",
		message,
		predict.RecommendedOrder,
		predict.ConfidenceScore*100,
	)

	return nil
}
