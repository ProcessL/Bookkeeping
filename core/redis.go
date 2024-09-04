package core

import (
	"context"
	"github.com/dotdancer/gogofly/global"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

var rdClient *redis.Client
var nDuration = 30 * 24 * 60 * 60 * time.Second

type RedisClient struct {
}

func InitRedis() {
	var client redis.UniversalClient
	client = redis.NewClient(&redis.Options{
		Addr:     global.Config.Redis.Addr,
		Password: global.Config.Redis.Password,
		DB:       global.Config.Redis.Db,
	})
	result, err := client.Ping(context.Background()).Result()
	if err != nil {
		global.Logger.Error("redis connect ping fail:", zap.Error(err))
		panic(err)
	} else {
		global.Logger.Info("redis connect ping result:", zap.String("result", result))
		global.Redis = client
	}
}

func (*RedisClient) Set(key string, value interface{}) error {
	return rdClient.Set(context.Background(), key, value, nDuration).Err()
}

func (*RedisClient) Get(key string) (any, error) {
	return rdClient.Get(context.Background(), key).Result()
}

func (*RedisClient) Delete(key ...string) error {
	return rdClient.Del(context.Background(), key...).Err()
}
