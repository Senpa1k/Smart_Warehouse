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

	// getting the key to access the AI
	apikey, err := config.Get("API_KEY")
	if err != nil {
		log.Printf("Failed to get Giga Api key: %v", err)
		return nil, err
	}

	// creating a client for communication with AI
	client, err := gigago.NewClient(ctx, apikey, gigago.WithCustomInsecureSkipVerify(true))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// getting data for analysis
	products, err := ai.repo.AIRequest(rq)
	if err != nil {
		return nil, err
	}

	// converting data to json format for further analysis
	assistantRequest, err := json.Marshal(products)
	if err != nil {
		return nil, err
	}

	// configuring the AI model
	model := client.GenerativeModel("GigaChat")
	model.SystemInstruction = "Ты - AI ассистент для анализа складских запасов. Анализируй данные инвентаризации и прогнозируй остатки товаров. Отвечай ТОЛЬКО в формате JSON."
	model.Temperature = 0.2
	model.TopP = 0.2
	model.MaxTokens = 3500
	model.RepetitionPenalty = 1.2

	// prompt for getting a forecast
	messages := []gigago.Message{
		{Role: gigago.RoleUser, Content: `Анализ складских запасов - прогноз на ` + strconv.Itoa(rq.PeriodDays) + ` дней. ДАННЫЕ ДЛЯ АНАЛИЗА:` + string(assistantRequest) +
			`

			ЗАДАЧА:
			Проанализируй тенденции потребления для каждого товара и спрогнозируй:
				1. Через сколько дней закончатся запасы (days_until_stockout)
				2. Рекомендуемое количество для заказа (recommended_order_quantity)
				3. Достоверность прогноза (confidence) от 0.0 до 1.0


			ТРЕБОВАНИЯ К ОТВЕТУ:
			- prediction_date должен быть: сегодняшней датой
			- Используй product_id и product_name из предоставленных данных
			- Ответ должен быть в точном JSON формате

			ФОРМАТ ОТВЕТА:
			{
				"predictions": [
					{
						"product_id": "string",
						"product_name: "string",
						"prediction_date": "dd.mm.yyyy",
						"days_until_stockout": "int",
						"recommended_order": "int",
						"confidence_score": "float",
					}
				],
				"confidence": "float",
			}

				Только JSON, без дополнительного текста.`,
		},
	}

	// request to the AI
	resp, err := model.Generate(ctx, messages)
	if err != nil {
		return nil, err
	}

	// bringing the data received from AI to a format convenient for further processing
	var aiResponse entities.AIResponse
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &aiResponse); err != nil {
		return nil, err
	}

	// writing the result to the database
	if err := ai.repo.AIResponse(aiResponse); err != nil {
		return nil, err
	}

	ai.made <- aiResponse

	return &aiResponse, nil
}
