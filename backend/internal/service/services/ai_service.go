package services

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Role1776/gigago"
	"github.com/Senpa1k/Smart_Warehouse/internal/config"
	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/repository"
	"github.com/sirupsen/logrus"
)

type AIService struct {
	repo  repository.AI
	redis repository.Redis
	made  chan<- interface{}
}

func NewAIService(repo repository.AI, made chan<- interface{}, redis repository.Redis) *AIService {
	return &AIService{
		repo:  repo,
		redis: redis,
		made:  made,
	}
}

func (ai *AIService) Predict(rq entities.AIRequest) (*entities.AIResponse, error) {
	// 1. Создаем ключ кеша на основе входных параметров
	cacheKey := fmt.Sprintf("ai:predict:%s:%d", generateRequestHash(rq), rq.PeriodDays)

	// 2. Пробуем получить из кеша
	if ai.redis != nil {
		if cached, err := ai.redis.Get(cacheKey); err == nil {
			var response entities.AIResponse
			if err := json.Unmarshal([]byte(cached), &response); err == nil {
				logrus.Infof("AI prediction served from cache for key: %s", cacheKey)
				return &response, nil
			}
		}
	}

	// 3. Если нет в кеше - выполняем AI запрос
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

	// 4. Сохраняем результат в кеш на 1 час
	if ai.redis != nil {
		data, _ := json.Marshal(aiResponse)
		ai.redis.Set(cacheKey, data, time.Hour)
		logrus.Infof("AI prediction cached for key: %s", cacheKey)
	}

	logrus.Infof("AI prediction cached for key: %s", cacheKey)

	ai.made <- aiResponse

	return &aiResponse, nil
}

// Вспомогательная функция для создания хеша запроса
func generateRequestHash(rq entities.AIRequest) string {
	// Создаем уникальный хеш на основе данных запроса
	data := fmt.Sprintf("%v", rq)
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}
