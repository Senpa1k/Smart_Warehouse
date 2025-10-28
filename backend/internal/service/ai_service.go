package service

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/Senpa1k/Smart_Warehouse/internal/config"
	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/repository"
	"github.com/paulrzcz/go-gigachat"
)




type AIService struct {
	repo repository.AI
	client *gigachat.Client
}


func NewAIService(repo repository.AI) *AIService {
	clientID, err1 := config.Get("GIGACHAT_CLIENT_ID")
	clientSecret, err2 := config.Get("GIGACHAT_CLIENT_SECRET")
	if err1 != nil {
		log.Printf("Failed to get GigaChat credentials: %v, %v", err1, err2)
		return nil
	}
	if err2 != nil {
		return nil
	}

	client, err := gigachat.NewClient(clientID, clientSecret)
	if err != nil {
		log.Printf("Failed to create GigaChat client: %v", err)
		return nil
	}

	return &AIService{repo: repo, client: client}
}

func (ai *AIService) Predict(rq entities.AIRequest) (*entities.AIResponse, error) {
	products, err := ai.repo.AIRequest(rq)
	if err != nil{
		return nil, err
	}
	
	assistantRequest, err := json.Marshal(products)
	if err != nil{
		return nil, err
	}


	resp, err := ai.client.Chat(&gigachat.ChatRequest{
		Model: gigachat.GIGACHAT_2_LITE,
		Messages: []gigachat.Message{
			{
				Role: gigachat.SystemRole,
				Content: "Ты - AI ассистент для анализа складских запасов. Анализируй данные инвентаризации и прогнозируй остатки товаров. Отвечай ТОЛЬКО в формате JSON.",
			},
			{
				Role: gigachat.AssistantRole,
				Content: `Проанализируй данные складских запасов указанные в этом json` + string(assistantRequest) + `и спрогнозируй остатки на количество дней равное ` + strconv.Itoa(rq.PeriodDays) +
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
		},
		Temperature:       	float64Ptr(0.2),
		TopP:             	float64Ptr(0.2),
		MaxTokens:			int64Ptr(2048),
		RepetitionPenalty:	float64Ptr(1.2),
		Stream:				boolPtr(false),  

	})

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
	
	return &aiResponse, nil
}

func float64Ptr(f float64) *float64 { return &f }
func int64Ptr(i int64) *int64       { return &i }  
func boolPtr(b bool) *bool          { return &b }
