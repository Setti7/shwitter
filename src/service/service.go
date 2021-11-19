package service

import (
	"github.com/Setti7/shwitter/session"
	"github.com/bsm/redislock"
	"github.com/gocql/gocql"
)

var services struct {
	Cassandra *gocql.Session
	Lock      *redislock.Client
	Session   *session.Session
}

// Initialize services
func Init() {
	Cassandra()
	Lock()
	Session()
}

func CleanUp() {
	Cassandra().Close()
}
