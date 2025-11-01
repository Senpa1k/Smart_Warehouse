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

	// 3. Если нет в кеше - выполняем AI запрос через HTTP API
	apikey, err := config.Get("API_KEY")
	if err != nil {
		log.Printf("Failed to get Giga Api key: %v", err)
		return nil, err
	}

	products, err := ai.repo.AIRequest(rq)
	if err != nil {
		return nil, err
	}

	assistantRequest, err := json.Marshal(products)
	if err != nil {
		return nil, err
	}

	// Подготовка запроса к GigaChat API
	messages := []GigaChatMessage{
		{
			Role:    "system",
			Content: "Ты - AI ассистент для анализа складских запасов. Анализируй данные инвентаризации и прогнозируй остатки товаров. Отвечай ТОЛЬКО в формате JSON.",
		},
		{
			Role: "user",
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
	}

	gigaRequest := GigaChatRequest{
		Model:       "GigaChat",
		Messages:    messages,
		Temperature: 0.2,
		TopP:        0.2,
		MaxTokens:   2000,
	}

	requestBody, err := json.Marshal(gigaRequest)
	if err != nil {
		return nil, err
	}

	// Создаем HTTP запрос
	req, err := http.NewRequest("POST", "https://gigachat.devices.sberbank.ru/api/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apikey)
	req.Header.Set("Accept", "application/json")

	// Выполняем запрос
	resp, err := ai.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GigaChat API error: %s, body: %s", resp.Status, string(body))
	}

	var gigaResponse GigaChatResponse
	if err := json.Unmarshal(body, &gigaResponse); err != nil {
		return nil, err
	}

	if len(gigaResponse.Choices) == 0 {
		return nil, fmt.Errorf("empty response from GigaChat API")
	}

	// Парсим AI ответ
	var aiResponse entities.AIResponse
	if err := json.Unmarshal([]byte(gigaResponse.Choices[0].Message.Content), &aiResponse); err != nil {
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

	ai.made <- aiResponse

	return &aiResponse, nil
}

// Вспомогательная функция для создания хеша запроса
func generateRequestHash(rq entities.AIRequest) string {
	data := fmt.Sprintf("%v", rq)
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}
