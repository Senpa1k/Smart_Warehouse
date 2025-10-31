package main

import (
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/repository"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	logrus.Info("Starting Redis connection test...")

	// Инициализация Redis
	redisClient, err := repository.NewRedisClient("redis://localhost:6379")
	if err != nil {
		logrus.Fatalf("Redis connection failed: %v", err)
	}
	defer redisClient.Close()

	logrus.Info("Redis connected successfully!")

	// Тест SET
	err = redisClient.Set("redis_test_key", "Hello Redis!", 30*time.Second)
	if err != nil {
		logrus.Fatalf("Redis SET failed: %v", err)
	}
	logrus.Info("Redis SET: OK")

	// Тест GET
	value, err := redisClient.Get("redis_test_key")
	if err != nil {
		logrus.Fatalf("Redis GET failed: %v", err)
	}
	logrus.Infof("Redis GET: %s", value)

	// Тест EXISTS
	exists, err := redisClient.Exists("redis_test_key")
	if err != nil {
		logrus.Fatalf("Redis EXISTS failed: %v", err)
	}
	logrus.Infof("Redis EXISTS: %v", exists)

	logrus.Info("All Redis tests passed! ✅")
}