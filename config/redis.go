package config

import (
	"os"
	"strconv"

	redisprom "github.com/globocom/go-redis-prometheus"
	"github.com/go-redis/redis/v8"
)

var Redis *redis.Client

func ConnectRedis() *redis.Client {
	hook := redisprom.NewHook(
		redisprom.WithInstanceName("cache"),
		redisprom.WithNamespace("blog"),
		redisprom.WithDurationBuckets([]float64{.001, .005, .01}),
	)

	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	redisUri := os.Getenv("REDIS_HOST") + ":6379"
	client := redis.NewClient(&redis.Options{
		Addr:     redisUri,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redisDB,
	})
	client.AddHook(hook)
	return client
}
