package service

import (
	"github.com/Setti7/shwitter/internal/config"
	"github.com/bsm/redislock"
	"github.com/gocql/gocql"
)

var conf *config.Config

var services struct {
	cassandra *gocql.Session
	lock      *redislock.Client
}

func SetConfig(c *config.Config) {
	if c == nil {
		panic("config is nil")
	}

	conf = c
}

// Initialize services
func Init() {
	Cassandra()
	Lock()
}

func CleanUp() {
	Cassandra().Close()
}
