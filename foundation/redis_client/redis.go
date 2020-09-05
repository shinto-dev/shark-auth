package redis_client

import (
	"fmt"

	"github.com/go-redis/redis/v7"
	"github.com/pkg/errors"
)

type Config struct {
	Host string
	Port int
}

func NewRedisClient(config Config) *redis.Client {
	dsn := fmt.Sprintf("%s:%d", config.Host, config.Port)

	client := redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(errors.Wrap(err, "error while connecting to redis"))
	}
	return client
}
