package service

import (
	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
	"sync"
)

var onceLock sync.Once

func initLock() {
	c := conf.Lock()

	client := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    c.Host,
	})
	services.lock = redislock.New(client)
}

func Lock() *redislock.Client {
	onceLock.Do(initLock)

	return services.lock
}
