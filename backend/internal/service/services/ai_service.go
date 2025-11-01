package services

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/config"
	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/repository"
	"github.com/sirupsen/logrus"
)

type AIService struct {
	repo       repository.AI
	redis      repository.Redis
	made       chan<- interface{}
	httpClient *http.Client
}

func NewAIService(repo repository.AI, made chan<- interface{}, redis repository.Redis) *AIService {
	// Создаем HTTP клиент с таймаутами и пропуском проверки SSL (аналогично WithCustomInsecureSkipVerify)
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	return &AIService{
		repo:       repo,
		redis:      redis,
		made:       made,
		httpClient: httpClient,
	}
}

// Структуры для GigaChat API
type GigaChatRequest struct {
	Model       string            `json:"model"`
	Messages    []GigaChatMessage `json:"messages"`
	Temperature float64           `json:"temperature,omitempty"`
	TopP        float64           `json:"top_p,omitempty"`
	MaxTokens   int               `json:"max_tokens,omitempty"`
}

type GigaChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GigaChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
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
	if err := json.Unmarshal([]byte(gigaResponse.Choices[0].Message.Content), &aiResponse); err != nil {
		return nil, err
	}

	// writing the result to the database
	if err := ai.repo.AIResponse(aiResponse); err != nil {
		return nil, err
	}

	// 4. Сохраняем результат в кеш на 1 час
	if ai.redis != nil {
		data, _ := json.Marshal(aiResponse)
		ai.redis.Set(cacheKey, data, time.Hour)
		logrus.Infof("AI prediction cached for key: %s", cacheKey)
	}

	ai.made <- aiResponse

	return &aiResponse, nil
}

// Вспомогательная функция для создания хеша запроса
func generateRequestHash(rq entities.AIRequest) string {
	data := fmt.Sprintf("%v", rq)
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}
