package redisclient

import (
	"fmt"

	"github.com/go-redis/redis/v7"
)

type Config struct {
	Host string
	Port int
}

// creates and returns a redis client with the configs passed.
func NewRedisClient(config Config) (*redis.Client, error) {
	dsn := fmt.Sprintf("%s:%d", config.Host, config.Port)

	client := redis.NewClient(&redis.Options{
		Addr: dsn,
	})
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
