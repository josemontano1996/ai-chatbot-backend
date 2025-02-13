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

func (r *Redis) Set(c *gin.Context, key string, value any, expiration time.Duration) *redis.StatusCmd {
	return r.client.Set(c, key, value, expiration)
}

func (r *Redis) Get(c *gin.Context, key string) *redis.StringCmd {
	return r.client.Get(c, key)
}

func (r *Redis) Delete(c *gin.Context, key string) *redis.IntCmd {
	return r.client.Del(c, key)
}

func (r *Redis) AddToList(c *gin.Context, key string, values ...any) (int64, error) {
	return r.client.RPush(c, key, values...).Result()
}

func (r *Redis) GetList(c *gin.Context, key string, start int64, stop int64) ([]string, error) {
	return r.client.LRange(c, key, start, stop).Result()
}
