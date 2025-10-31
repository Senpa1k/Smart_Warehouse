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

func (ai *AIPostgres) AIRequest(rq entities.AIRequest) (*[]models.Products, error) {
	var products []models.Products
	err := ai.db.Where("category IN ?", rq.Categories).Find(&products).Error
	if err != nil {
		return nil, err
	}

	return &products, nil
}

func (ai *AIPostgres) AIResponse(rp entities.AIResponse) error {
	for _, elem := range rp.Predictions {
		predictionDate := time.Now()

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
