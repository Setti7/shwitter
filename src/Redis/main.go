package Redis

import (
	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
)

var Locker *redislock.Client

func ConnectToRedis() {
	client := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    "127.0.0.1:6379",
	})
	Locker = redislock.New(client)
}
