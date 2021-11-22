package config

import "github.com/urfave/cli/v2"

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
	c := &LockConfig{}

	if ctx == nil {
		return &LockDefault
	}

	c.Host = ctx.String("redis-host")

	return c
}
