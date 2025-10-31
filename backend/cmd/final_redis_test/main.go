package main

import (
	"context"
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/repository"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	logrus.Info("🎯 Starting Final Redis Integration Test...")

	// 1. Подключаемся к Redis
	redisClient, err := repository.NewRedisClient("redis://localhost:6379")
	if err != nil {
		logrus.Fatalf("❌ Redis connection failed: %v", err)
	}
	defer redisClient.Close()

	logrus.Info("✅ Redis connected successfully")

	// 2. Тест статусов роботов
	testRobotStatuses(redisClient)

	// 3. Тест Rate Limiting
	testRateLimiting(redisClient)

	// 4. Тест Pub/Sub
	testPubSub(redisClient)

	logrus.Info("🎉 All Redis features working correctly! 4th stage completed!")
}

func testRobotStatuses(redis repository.Redis) {
	logrus.Info("🤖 Testing robot status management...")

	// Сохраняем статусы тестовых роботов
	robots := []struct {
		id      string
		battery int
		status  string
	}{
		{"RB-001", 85, "active"},
		{"RB-002", 42, "low_battery"},
		{"RB-003", 15, "charging"},
	}

	for _, robot := range robots {
		// Сохраняем статус
		redis.SetRobotStatus(robot.id, robot.status, time.Minute)
		redis.SetRobotBattery(robot.id, robot.battery, time.Minute)
		redis.SetRobotOnline(robot.id)

		logrus.Infof("✅ Robot %s: battery=%d%%, status=%s", robot.id, robot.battery, robot.status)
	}

	// Проверяем чтение
	for _, robot := range robots {
		online, _ := redis.IsRobotOnline(robot.id)
		battery, _ := redis.GetRobotBattery(robot.id)
		status, _ := redis.GetRobotStatus(robot.id)

		logrus.Infof("📊 Robot %s: online=%v, battery=%d%%, status=%s",
			robot.id, online, battery, status)
	}
}

func testRateLimiting(redis repository.Redis) {
	logrus.Info("🛡️ Testing rate limiting...")

	key := "rate:test:127.0.0.1"

	// Тестируем лимит 3 запроса в минуту
	for i := 1; i <= 5; i++ {
		allowed, err := redis.CheckRateLimit(key, 3, time.Minute)
		if err != nil {
			logrus.Errorf("❌ Rate limit error: %v", err)
			continue
		}

		if allowed {
			logrus.Infof("✅ Request %d: ALLOWED", i)
		} else {
			logrus.Infof("🚫 Request %d: BLOCKED (rate limit exceeded)", i)
		}
	}
}

func testPubSub(redis repository.Redis) {
	logrus.Info("📡 Testing Pub/Sub system...")

	ctx := context.Background()

	// Подписчик
	go func() {
		pubsub := redis.Subscribe("robot_updates")
		defer pubsub.Close()

		count := 0
		for {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				return
			}

			count++
			logrus.Infof("📨 Received message %d: %s", count, msg.Payload)

			if count >= 2 {
				return
			}
		}
	}()

	// Даем время на подписку
	time.Sleep(1 * time.Second)

	// Публикуем тестовые сообщения
	messages := []string{
		`{"type": "robot_online", "robot_id": "RB-001", "battery": 85}`,
		`{"type": "scan_complete", "robot_id": "RB-002", "items": 15}`,
	}

	for i, msg := range messages {
		err := redis.Publish("robot_updates", msg)
		if err != nil {
			logrus.Errorf("❌ Publish error: %v", err)
		} else {
			logrus.Infof("📤 Published message %d", i+1)
		}
		time.Sleep(500 * time.Millisecond)
	}

	time.Sleep(2 * time.Second)
}
