package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

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

	productsWithHistory, err := ai.repo.AIRequest(rq)
	if err != nil {
		return nil, err
	}

	if len(*productsWithHistory) == 0 {
		return &entities.AIResponse{
			Predictions: []entities.Predictions{},
			Confidence:  0,
		}, nil
	}

	assistantRequest, err := json.Marshal(productsWithHistory)
	if err != nil {
		return nil, err
	}

	model := client.GenerativeModel("GigaChat")
	model.SystemInstruction = "Ты - AI ассистент для анализа складских запасов. Анализируй данные инвентаризации и прогнозируй остатки товаров. Отвечай ТОЛЬКО в формате JSON."
	model.Temperature = 0.3
	model.TopP = 0.3
	model.MaxTokens = 2048
	model.RepetitionPenalty = 1.1

	todayStr := time.Now().Format("02.01.2006")
	
	messages := []gigago.Message{
		{Role: gigago.RoleUser, Content: fmt.Sprintf(`Проанализируй данные складских запасов:

%s

Сегодняшняя дата: %s
Период прогноза: %d дней

Для каждого товара:
1. Проанализируй историю сканирований (recent_scans) за период
2. Учти текущий остаток (current_stock)
3. Учти среднее ежедневное потребление (average_daily_consumption)
4. Учти минимальный запас (min_stock) и оптимальный запас (optimal_stock)
5. Спрогнозируй:
   - Через сколько дней закончатся запасы (days_until_stockout)
   - Рекомендуемое количество для заказа (recommended_order) - должно быть >= optimal_stock
   - Достоверность прогноза (confidence_score) от 0.0 до 1.0

Верни ответ СТРОГО в формате JSON (prediction_date в формате dd.mm.yyyy):
{
  "predictions": [
    {
      "product_id": "TEL-4567",
      "prediction_date": "%s",
      "days_until_stockout": 15,
      "recommended_order": 100,
      "confidence_score": 0.85
    }
  ],
  "confidence": 0.85
}

ВАЖНО: Верни ТОЛЬКО JSON, без дополнительного текста, markdown или объяснений.`, 
			string(assistantRequest), todayStr, rq.PeriodDays, todayStr)},
	}

	resp, err := model.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("GigaChat API error: %w", err)
	}

	content := resp.Choices[0].Message.Content
	
	// Try to extract JSON if wrapped in markdown
	if len(content) > 0 {
		// Remove markdown code blocks if present
		if content[0] == '`' {
			start := 0
			end := len(content)
			
			// Find first {
			for i, c := range content {
				if c == '{' {
					start = i
					break
				}
			}
			
			// Find last }
			for i := len(content) - 1; i >= 0; i-- {
				if content[i] == '}' {
					end = i + 1
					break
				}
			}
			
			content = content[start:end]
		}
	}

	var aiResponse entities.AIResponse
	if err := json.Unmarshal([]byte(content), &aiResponse); err != nil {
		log.Printf("Failed to parse AI response. Raw content: %s", content)
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	// Save predictions to database
	if err := ai.repo.AIResponse(aiResponse); err != nil {
		log.Printf("Failed to save AI predictions: %v", err)
		// Don't return error, just log it
	}

	if ai.made != nil {
		select {
		case ai.made <- aiResponse:
		default:
			// Channel full or closed, skip
		}
	}

	return &aiResponse, nil
}
