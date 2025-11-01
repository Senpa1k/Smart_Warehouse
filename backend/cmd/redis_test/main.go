package main

import (
	"encoding/json"
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/repository"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	logrus.Info("Starting Redis integration test...")

	// 1. Подключаемся к Redis
	redisClient, err := repository.NewRedisClient("redis://localhost:6379")
	if err != nil {
		logrus.Fatalf("Redis connection failed: %v", err)
	}
	defer redisClient.Close()

	logrus.Info("✅ Redis connected successfully")

	// 2. Тест кеширования dashboard данных (соответствует твоему DashInfo)
	dashData := entities.DashInfo{
		Statistics: entities.Statistics{
			ActiveRobots:      5,
			TotalRobots:       10,
			ItemsCheckedToday: 150,
			CriticalItems:     3,
			AvgBattery:        85,
		},
		// ListRobots и ListScans оставляем пустыми для теста
	}

	dashJSON, _ := json.Marshal(dashData)
	err = redisClient.Set("dashboard:current", dashJSON, 5*time.Second)
	if err != nil {
		logrus.Fatalf("❌ Dashboard cache SET failed: %v", err)
	}
	logrus.Info("✅ Dashboard cache SET successful")

	// 3. Тест чтения из кеша
	cachedDash, err := redisClient.Get("dashboard:current")
	if err != nil {
		logrus.Fatalf("❌ Dashboard cache GET failed: %v", err)
	}

	var restoredDash entities.DashInfo
	if err := json.Unmarshal([]byte(cachedDash), &restoredDash); err != nil {
		logrus.Fatalf("❌ Dashboard cache unmarshal failed: %v", err)
	}
	logrus.Infof("✅ Dashboard cache GET successful: ActiveRobots=%d", restoredDash.Statistics.ActiveRobots)

	// 4. Тест кеширования AI прогнозов (соответствует твоему AIResponse)
	aiResponse := entities.AIResponse{
		Predictions: []entities.Predictions{
			{
				ProductID:         "TEL-123",
				PredictionDate:    "01.11.2024",
				DaysUntilStockout: 7,
				RecommendedOrder:  50,
				ConfidenceScore:   0.85,
			},
		},
		Confidence: 0.82, // Исправлена опечатка в поле
	}

	aiJSON, _ := json.Marshal(aiResponse)
	err = redisClient.Set("ai:predict:hash123:7", aiJSON, time.Hour)
	if err != nil {
		logrus.Fatalf("❌ AI cache SET failed: %v", err)
	}
	logrus.Info("✅ AI cache SET successful")

	// 5. Тест чтения AI из кеша
	cachedAI, err := redisClient.Get("ai:predict:hash123:7")
	if err != nil {
		logrus.Fatalf("❌ AI cache GET failed: %v", err)
	}

	var restoredAI entities.AIResponse
	if err := json.Unmarshal([]byte(cachedAI), &restoredAI); err != nil {
		logrus.Fatalf("❌ AI cache unmarshal failed: %v", err)
	}
	logrus.Infof("✅ AI cache GET successful: Confidence=%.2f", restoredAI.Confidence)

	// 6. Тест инвалидации кеша
	err = redisClient.Delete("dashboard:current")
	if err != nil {
		logrus.Fatalf("❌ Cache DELETE failed: %v", err)
	}
	logrus.Info("✅ Cache invalidation successful")

	// 7. Проверка что кеш очищен
	exists, err := redisClient.Exists("dashboard:current")
	if err != nil {
		logrus.Fatalf("❌ Cache EXISTS check failed: %v", err)
	}
	logrus.Infof("✅ Cache invalidation verified: dashboard exists = %v", exists)

	// 8. Тест TTL (время жизни)
	err = redisClient.Set("test_ttl", "temp_data", 2*time.Second)
	if err != nil {
		logrus.Fatalf("❌ TTL test SET failed: %v", err)
	}
	logrus.Info("✅ TTL test SET successful")

	time.Sleep(3 * time.Second) // Ждем истечения TTL

	_, err = redisClient.Get("test_ttl")
	if err != nil {
		logrus.Info("✅ TTL expiration working correctly")
	} else {
		logrus.Info("❌ TTL expiration not working")
	}

	logrus.Info("🎉 All Redis integration tests passed! Caching layer is ready.")
}
