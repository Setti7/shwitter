package service

import (
	"github.com/Setti7/shwitter/internal/config"
	"github.com/gocql/gocql"
	"sync"
)

var onceCassandra sync.Once

func initCassandra() {
	cluster := gocql.NewCluster(config.CassandraDefault.Hosts...)
	cluster.Keyspace = config.CassandraDefault.Keyspace

	sess, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}

	services.Cassandra = sess
}

func Cassandra() *gocql.Session {
	onceCassandra.Do(initCassandra)

	return services.Cassandra
}
