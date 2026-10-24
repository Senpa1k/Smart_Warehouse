// Package repository defines persistence contracts and Redis operations.
package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient wraps Redis caching, pub/sub, robot state, and rate limiting.
type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisClient connects to Redis using a URL and verifies the connection.
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

// Set stores a value with an expiration time.
func (r *RedisClient) Set(key string, value any, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

// Get reads a cached string value.
func (r *RedisClient) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

// Delete removes a cached key.
func (r *RedisClient) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

// Exists reports whether a cached key is present.
func (r *RedisClient) Exists(key string) (bool, error) {
	result, err := r.client.Exists(r.ctx, key).Result()
	return result > 0, err
}

// Close closes the underlying Redis connection.
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// HSet writes fields into a Redis hash.
func (r *RedisClient) HSet(key string, values ...any) error {
	return r.client.HSet(r.ctx, key, values...).Err()
}

// HGet reads a field from a Redis hash.
func (r *RedisClient) HGet(key, field string) (string, error) {
	return r.client.HGet(r.ctx, key, field).Result()
}

// Publish - отправка сообщения в Redis channel
func (r *RedisClient) Publish(channel string, message any) error {
	return r.client.Publish(r.ctx, channel, message).Err()
}

// Subscribe - подписка на Redis channel
func (r *RedisClient) Subscribe(channel string) *redis.PubSub {
	return r.client.Subscribe(r.ctx, channel)
}

// SetRobotStatus caches a robot status for the given duration.
func (r *RedisClient) SetRobotStatus(robotID, status string, expiration time.Duration) error {
	key := fmt.Sprintf("robot:%s:status", robotID)
	return r.Set(key, status, expiration)
}

// GetRobotStatus reads the cached status of a robot.
func (r *RedisClient) GetRobotStatus(robotID string) (string, error) {
	key := fmt.Sprintf("robot:%s:status", robotID)
	return r.Get(key)
}

// SetRobotBattery caches a robot battery level for the given duration.
func (r *RedisClient) SetRobotBattery(robotID string, batteryLevel int, expiration time.Duration) error {
	key := fmt.Sprintf("robot:%s:battery", robotID)
	return r.Set(key, strconv.Itoa(batteryLevel), expiration)
}

// GetRobotBattery reads the cached battery level of a robot.
func (r *RedisClient) GetRobotBattery(robotID string) (int, error) {
	key := fmt.Sprintf("robot:%s:battery", robotID)
	val, err := r.Get(key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(val)
}

// SetRobotOnline records the latest time a robot was seen.
func (r *RedisClient) SetRobotOnline(robotID string) error {
	key := fmt.Sprintf("robot:%s:last_seen", robotID)
	return r.Set(key, time.Now().Format(time.RFC3339), 30*time.Second)
}

// IsRobotOnline reports whether a robot was seen recently.
func (r *RedisClient) IsRobotOnline(robotID string) (bool, error) {
	key := fmt.Sprintf("robot:%s:last_seen", robotID)
	return r.Exists(key)
}

// CheckRateLimit increments a request counter and checks its window limit.
func (r *RedisClient) CheckRateLimit(key string, limit int, window time.Duration) (bool, error) {
	count, err := r.client.Incr(r.ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("increment rate limit counter: %w", err)
	}
	if count == 1 {
		if err := r.client.Expire(r.ctx, key, window).Err(); err != nil {
			return false, fmt.Errorf("expire rate limit counter: %w", err)
		}
	}
	return count <= int64(limit), nil
}
