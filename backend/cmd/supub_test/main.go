package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/repository"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	logrus.Info("Starting Redis Pub/Sub test...")

	// Подключаемся к Redis
	redisClient, err := repository.NewRedisClient("redis://localhost:6379")
	if err != nil {
		logrus.Fatalf("Redis connection failed: %v", err)
	}
	defer redisClient.Close()

	// Тестируем Pub/Sub
	testRedisPubSub()
}

func testRedisPubSub() {
	ctx := context.Background()

	// Создаем отдельного клиента для теста
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	// 1. Подписываемся на канал
	go func() {
		pubsub := client.Subscribe(ctx, "robot_updates")
		defer pubsub.Close()

		logrus.Info("📡 Subscribed to robot_updates channel")

		for {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				logrus.Errorf("Receive message error: %v", err)
				return
			}

			logrus.Infof("📨 Received message: %s", msg.Payload)
		}
	}()

	// 2. Даем время на подписку
	time.Sleep(1 * time.Second)

	// 3. Публикуем тестовые сообщения
	for i := 1; i <= 3; i++ {
		robotData := map[string]interface{}{
			"robot_id":      fmt.Sprintf("RB-%03d", i),
			"battery_level": 80 + i*5,
			"timestamp":     time.Now().Format(time.RFC3339),
			"message_type":  "robot_update",
		}

		data, _ := json.Marshal(robotData)

		err := client.Publish(ctx, "robot_updates", string(data)).Err()
		if err != nil {
			logrus.Errorf("Publish error: %v", err)
		} else {
			logrus.Infof("📤 Published message %d", i)
		}

		time.Sleep(2 * time.Second)
	}

	logrus.Info("✅ Redis Pub/Sub test completed!")
}
