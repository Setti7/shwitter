package config

import (
	"github.com/urfave/cli/v2"
)

type Config struct {
	cassandra *CassandraConfig
	redis     *RedisConfig
}

func NewConfig(ctx *cli.Context) *Config {
	c := &Config{
		cassandra: NewCassandraConfig(ctx),
		redis:     NewRedisConfig(ctx),
	}

	return c
}
