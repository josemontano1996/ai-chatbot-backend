package repository

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(addr string, password string, db int) *Redis {
	return &Redis{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}),
	}

}

func (r *Redis) Set(c *gin.Context, key string, value any, expiration time.Duration) error {
	return r.client.Set(c, key, value, expiration).Err()
}

func (r *Redis) Get(c *gin.Context, key string) (val any, err error) {
	return r.client.Get(c, key).Result()
}

