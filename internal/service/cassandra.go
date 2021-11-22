package service

import (
	"github.com/gocql/gocql"
	"sync"
)

var onceCassandra sync.Once

func initCassandra() {
	c := conf.Cassandra()

	cluster := gocql.NewCluster(c.Hosts...)
	cluster.Keyspace = c.Keyspace

	sess, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}

	services.cassandra = sess
}

func Cassandra() *gocql.Session {
	onceCassandra.Do(initCassandra)

	return services.cassandra
}
