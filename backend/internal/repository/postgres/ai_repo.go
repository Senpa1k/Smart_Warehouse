package postgres

import (
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/models"
	"gorm.io/gorm"
)

type AIPostgres struct {
	db *gorm.DB
}

func NewAIPostgres(db *gorm.DB) *AIPostgres {
	return &AIPostgres{db: db}
}

// getting the necessary data to send a request to the AI
func (ai *AIPostgres) AIRequest(rq entities.AIRequest) ([]models.InventoryHistory, error) {
	var products []models.InventoryHistory

	startTime := time.Now().Add(-72 * time.Hour)

	// creating an additional query to split the data into intervals and reduce the total number
	subQuery := ai.db.
		Table("inventory_history").
		Select(`
			product_id,
			DATE_TRUNC('hour', scanned_at) + 
			INTERVAL '10 min' * FLOOR(EXTRACT(minute FROM scanned_at) / 10) as time_slot,
			MAX(scanned_at) as latest_in_slot`).
		Where("scanned_at >= ?", startTime).
		Group("product_id, time_slot")

	// create a query to get the necessary data from the database and select by category if any
	query := ai.db.Select(`
						inventory_history.id,
						inventory_history.product_id,
						inventory_history.quantity,
						inventory_history.status,
						inventory_history.scanned_at`).
					Preload("Product").Joins("JOIN products ON inventory_history.product_id = products.id").
					Joins("JOIN (?) as time_slots ON inventory_history.product_id = time_slots.product_id AND inventory_history.scanned_at = time_slots.latest_in_slot", subQuery)
	if len(rq.Categories) > 0 {
		query = query.Where("products.category IN ? ", rq.Categories)
	}

	err := query.Order("inventory_history.scanned_at DESC").Find(&products).Error
	if err != nil {
		return nil, err
	}


	return products, nil
}

// recording the response from the AI in the database
func (ai *AIPostgres) AIResponse(rp entities.AIResponse) error {
	for _, elem := range rp.Predictions {
		predictionDate, err := time.Parse("2006-01-02", elem.PredictionDate)
			if err != nil {
				predictionDate, err = time.Parse("02.01.2006", elem.PredictionDate)
				if err != nil {
					return err
				}
			}

		var prediction models.AiPrediction = models.AiPrediction{
			ProductID:         elem.ProductID,
			PredictionDate:    predictionDate,
			DaysUntilStockout: elem.DaysUntilStockout,
			RecommendedOrder:  elem.RecommendedOrder,
			ConfidenceScore:   elem.ConfidenceScore,
		}

		if err := ai.db.Create(&prediction).Error; err != nil {
			return err
		}
	}

	return nil
}
