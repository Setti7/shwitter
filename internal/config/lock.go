package config

import (
	"github.com/urfave/cli/v2"
	"os"
)

type LockConfig struct {
	Host string
}

var LockDefault = LockConfig{
	Host: "127.0.0.1:6379",
}

func (c *Config) Lock() *LockConfig {
	return c.lock
}

func NewLockConfig(ctx *cli.Context) *LockConfig {
	c := &LockDefault

	getLockConfigFromEnv(c)
	getLockConfigFromCLI(c, ctx)

	return c
}

func getLockConfigFromEnv(c *LockConfig) {
	c.Host = os.Getenv(REDIS_HOST_ENV)
}

func getLockConfigFromCLI(c *LockConfig, ctx *cli.Context) {
	c.Host = ctx.String(REDIS_HOST_FLAG_NAME)
}
