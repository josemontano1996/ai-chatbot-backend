package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/josemontano1996/ai-chatbot-backend/internal/entities"
	"github.com/redis/go-redis/v9"
)

type RedisMessageRepository struct {
	client *redis.Client
}

type redisConfig struct {
	Addr     string
	Password string
	DB       int
}

func NewRedisRepository(config *redisConfig) *RedisMessageRepository {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	return &RedisMessageRepository{
		client: client,
	}
}

func NewRedisConfig(address string, password string, db int) *redisConfig {
	return &redisConfig{
		Addr:     address,
		Password: password,
		DB:       db,
	}
}

func (r *RedisMessageRepository) GetChatHistory(ctx context.Context, key string) (*entities.ChatHistory, error) {
	jsonHistory, err := r.client.LRange(ctx, key, 0, -1).Result() // Get the entire list
	if err != nil {
		return nil, fmt.Errorf("error getting messages from Redis: %w", err)
	}

	prevHistory, err := entities.ParseArrayToChatHistory(jsonHistory)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON history to struct: %w", err)
	}
	return prevHistory, nil
}

func (r *RedisMessageRepository) SaveMessage(ctx context.Context, key string, userMsg *entities.ChatMessage) error {
	userMsgJSON, err := json.Marshal(userMsg)

	if err != nil {
		return fmt.Errorf("error marshaling user message to JSON: %w", err)
	}

	_, err = r.client.RPush(ctx, key, userMsgJSON).Result()

	if err != nil {
		return fmt.Errorf("error pushing responses to Redis: %w", err)
	}

	return nil
}

func (r *RedisMessageRepository) SaveMessages(ctx context.Context, key string, messages ...*entities.ChatMessage) error {
	jsonMessages := make([]interface{}, len(messages))
	for i, msg := range messages {
		msgJSON, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("error marshaling message to JSON: %w, message: %+v", err, msg)
		}
		jsonMessages[i] = msgJSON
	}

	if len(jsonMessages) > 0 {
		_, err := r.client.RPush(ctx, key, jsonMessages...).Result()
		if err != nil {
			return fmt.Errorf("error pushing messages to Redis: %w", err)
		}
	}

	return nil
}

func (r *RedisMessageRepository) Close() error {
	if r.client != nil {
		if err := r.client.Close(); err != nil {
			log.Printf("Error closing Redis client: %v", err)
			return err // Or just log, depending on error handling policy
		}
	}
	return nil
}
