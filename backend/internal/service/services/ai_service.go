package services

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/Role1776/gigago"
	"github.com/Senpa1k/Smart_Warehouse/internal/config"
	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/repository"
)

type AIService struct {
	repo repository.AI
	made chan<- interface{}
}

func NewAIService(repo repository.AI, made chan<- interface{}) *AIService {
	return &AIService{repo: repo, made: made}
}

func (ai *AIService) Predict(rq entities.AIRequest) (*entities.AIResponse, error) {
	ctx := context.Background()

	apikey, err := config.Get("API_KEY")
	if err != nil {
		log.Printf("Failed to get Giga Api key: %v", err)
		return nil, err
	}

	client, err := gigago.NewClient(ctx, apikey, gigago.WithCustomInsecureSkipVerify(true))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	products, err := ai.repo.AIRequest(rq)
	if err != nil {
		return nil, err
	}

	assistantRequest, err := json.Marshal(products)
	if err != nil {
		return nil, err
	}

	model := client.GenerativeModel("GigaChat")
	model.SystemInstruction = "Ты - AI ассистент для анализа складских запасов. Анализируй данные инвентаризации и прогнозируй остатки товаров. Отвечай ТОЛЬКО в формате JSON."
	model.Temperature = 0.2
	model.TopP = 0.2
	model.MaxTokens = 2000
	model.RepetitionPenalty = 1.2

	messages := []gigago.Message{
		{Role: gigago.RoleUser, Content: `Проанализируй данные складских запасов указанные в этом json` + string(assistantRequest) + `и спрогнозируй остатки на количество дней равное ` + strconv.Itoa(rq.PeriodDays) +
			` Проанализируй тенденции потребления для каждого товара и спрогнозируй:
				1. Через сколько дней закончатся запасы (days_until_stockout)
				2. Рекомендуемое количество для заказа (recommended_order_quantity)
				3. Достоверность прогноза (confidence) от 0.0 до 1.0

				Верни ответ в формате JSON, напиши prediction_date в формате dd.mm.yyyy:
				{
					"predictions": [
						{
							"product_id": string,
							"prediction_date": string
							"days_until_stockout": int,
							"recommended_order": int,
							"confidence_score": float,
						}
					]
					"confidence": float,
				}

				Только JSON, без дополнительного текста.`,
		},
	}

	resp, err := model.Generate(ctx, messages)
	if err != nil {
		return nil, err
	}

	var aiResponse entities.AIResponse
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &aiResponse); err != nil {
		return nil, err
	}

	if err := ai.repo.AIResponse(aiResponse); err != nil {
		return nil, err
	}

	ai.made <- aiResponse

	return &aiResponse, nil
}
