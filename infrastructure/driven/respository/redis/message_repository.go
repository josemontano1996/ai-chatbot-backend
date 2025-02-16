package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/josemontano1996/ai-chatbot-backend/internal/entities"
	"github.com/redis/go-redis/v9"
)

type RedisMessageRepository struct {
	client *redis.Client
	config *redisConfig
}

type redisConfig struct {
	Addr               string        `json:"addr" binding:"required"`
	Password           string        `json:"password" binding:"required"`
	DB                 int           `json:"db" binding:"required"`
	ExpirationDuration time.Duration `json:"expiration_duration" binding:"required"`
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

func NewRedisConfig(address string, password string, db int, expirationDuration time.Duration) *redisConfig {
	config := &redisConfig{
		Addr:               address,
		Password:           password,
		DB:                 db,
		ExpirationDuration: expirationDuration,
	}

	validate := validator.New()

	err := validate.Struct(config)

	if err != nil {
		log.Fatalf("Error validating Redis config: %v", err)
	}
	return config

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
	pipe := r.client.Pipeline()

	pipe.RPush(ctx, key, userMsgJSON)

	pipe.Expire(ctx, key, r.config.ExpirationDuration)

	_, err = pipe.Exec(ctx)

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
		pipe := r.client.Pipeline()

		pipe.RPush(ctx, key, jsonMessages...)

		pipe.Expire(ctx, key, r.config.ExpirationDuration)

		_, err := pipe.Exec(ctx)

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

