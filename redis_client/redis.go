package redis_client

import (
	"github.com/go-redis/redis/v7"
)

var  Client *redis.Client

func init() {
	dsn := "localhost:6379"

	Client = redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	_, err := Client.Ping().Result()
	if err != nil {
		panic(err)
	}
}
