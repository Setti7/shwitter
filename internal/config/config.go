package config

import (
	"github.com/urfave/cli/v2"
)

type Config struct {
	cassandra *CassandraConfig
	lock      *LockConfig // TODO change to RedisCondig
}

func NewConfig(ctx *cli.Context) *Config {
	c := &Config{
		cassandra: NewCassandraConfig(ctx),
		lock:      NewLockConfig(ctx),
	}

	return c
}
