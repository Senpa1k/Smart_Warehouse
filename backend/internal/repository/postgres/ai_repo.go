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

func (ai *AIPostgres) AIRequest(rq entities.AIRequest) (*[]entities.ProductWithHistory, error) {
	var products []models.Products
	
	// Get products by categories or all if empty
	query := ai.db.Model(&models.Products{})
	if len(rq.Categories) > 0 {
		query = query.Where("category IN ?", rq.Categories)
	}
	
	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}

	// Calculate date range for history
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -rq.PeriodDays)

	var result []entities.ProductWithHistory
	
	for _, product := range products {
		// Get recent scans for this product
		var history []models.InventoryHistory
		ai.db.Where("product_id = ? AND scanned_at >= ? AND scanned_at <= ?", 
			product.ID, startDate, endDate).
			Order("scanned_at DESC").
			Limit(50).
			Find(&history)
		
		// Calculate current stock (last scan)
		currentStock := 0
		if len(history) > 0 {
			currentStock = history[0].Quantity
		}
		
		// Calculate average daily consumption
		avgDaily := 0.0
		if len(history) > 1 {
			firstQty := history[len(history)-1].Quantity
			lastQty := history[0].Quantity
			days := float64(rq.PeriodDays)
			if days > 0 && firstQty > lastQty {
				avgDaily = float64(firstQty-lastQty) / days
			}
		}
		
		result = append(result, entities.ProductWithHistory{
			Product:      product,
			History:      history,
			CurrentStock: currentStock,
			AverageDaily: avgDaily,
		})
	}

	return &result, nil
}

func (ai *AIPostgres) AIResponse(rp entities.AIResponse) error {
	for _, elem := range rp.Predictions {
		predictionDate, err := time.Parse("02.01.2006", elem.PredictionDate)
		if err != nil {
			// Try alternative format
			predictionDate, err = time.Parse("2006-01-02", elem.PredictionDate)
			if err != nil {
				// Use current date if parsing fails
				predictionDate = time.Now()
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
