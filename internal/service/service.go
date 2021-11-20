package service

import (
	"github.com/bsm/redislock"
	"github.com/gocql/gocql"
)

var services struct {
	Cassandra *gocql.Session
	Lock      *redislock.Client
}

// Initialize services
func Init() {
	Cassandra()
	Lock()
}

func CleanUp() {
	Cassandra().Close()
}
