package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisClient(redisURL string) (*RedisClient, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opts)

	// Проверка подключения
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisClient{
		client: client,
		ctx:    ctx,
	}, nil
}

func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

func (r *RedisClient) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

func (r *RedisClient) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

func (r *RedisClient) Exists(key string) (bool, error) {
	result, err := r.client.Exists(r.ctx, key).Result()
	return result > 0, err
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

// Дополнительные методы для будущего использования
func (r *RedisClient) HSet(key string, values ...interface{}) error {
	return r.client.HSet(r.ctx, key, values...).Err()
}

func (r *RedisClient) HGet(key, field string) (string, error) {
	return r.client.HGet(r.ctx, key, field).Result()
}

// Publish - отправка сообщения в Redis channel
func (r *RedisClient) Publish(channel string, message interface{}) error {
	return r.client.Publish(r.ctx, channel, message).Err()
}

// Subscribe - подписка на Redis channel
func (r *RedisClient) Subscribe(channel string) *redis.PubSub {
	return r.client.Subscribe(r.ctx, channel)
}

// Robot status management
func (r *RedisClient) SetRobotStatus(robotID, status string, expiration time.Duration) error {
	key := fmt.Sprintf("robot:%s:status", robotID)
	return r.Set(key, status, expiration)
}

func (r *RedisClient) GetRobotStatus(robotID string) (string, error) {
	key := fmt.Sprintf("robot:%s:status", robotID)
	return r.Get(key)
}

func (r *RedisClient) SetRobotBattery(robotID string, batteryLevel int, expiration time.Duration) error {
	key := fmt.Sprintf("robot:%s:battery", robotID)
	return r.Set(key, strconv.Itoa(batteryLevel), expiration)
}

func (r *RedisClient) GetRobotBattery(robotID string) (int, error) {
	key := fmt.Sprintf("robot:%s:battery", robotID)
	val, err := r.Get(key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(val)
}

// Online robots tracking
func (r *RedisClient) SetRobotOnline(robotID string) error {
	key := fmt.Sprintf("robot:%s:last_seen", robotID)
	return r.Set(key, time.Now().Format(time.RFC3339), 30*time.Second)
}

func (r *RedisClient) IsRobotOnline(robotID string) (bool, error) {
	key := fmt.Sprintf("robot:%s:last_seen", robotID)
	return r.Exists(key)
}

// Rate limiting - защита от слишком частых запросов
func (r *RedisClient) CheckRateLimit(key string, limit int, window time.Duration) (bool, error) {
	// Получаем текущее значение
	current, err := r.Get(key)
	if err != nil {
		// Первый запрос - устанавливаем счетчик
		r.Set(key, "1", window)
		return true, nil
	}

	// Преобразуем в число
	count, err := strconv.Atoi(current)
	if err != nil {
		return false, err
	}

	// Проверяем лимит
	if count >= limit {
		return false, nil // Лимит исчерпан
	}

	// Увеличиваем счетчик
	r.Set(key, strconv.Itoa(count+1), window)
	return true, nil
}
