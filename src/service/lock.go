package service

import (
	"github.com/Setti7/shwitter/config"
	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
	"sync"
)

var onceLock sync.Once

func initLock() {
	client := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    config.LockDefault.Host + ":" + config.LockDefault.Port,
	})
	services.Lock = redislock.New(client)
}

func Lock() *redislock.Client {
	onceLock.Do(initLock)

	return services.Lock
}
