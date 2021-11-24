package config

import (
	"github.com/urfave/cli/v2"
	"os"
)

type RedisConfig struct {
	Host string
}

var RedisDefault = RedisConfig{
	Host: "127.0.0.1:6379",
}

func (c *Config) Redis() *RedisConfig {
	return c.redis
}

func NewRedisConfig(ctx *cli.Context) *RedisConfig {
	c := &RedisDefault

	getRedisConfigFromEnv(c)
	getRedisConfigFromCLI(c, ctx)

	return c
}

func getRedisConfigFromEnv(c *RedisConfig) {
	c.Host = os.Getenv(REDIS_HOST_ENV)
}

func getRedisConfigFromCLI(c *RedisConfig, ctx *cli.Context) {
	c.Host = ctx.String(REDIS_HOST_FLAG_NAME)
}
