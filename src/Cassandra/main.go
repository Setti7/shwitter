package Cassandra

import (
	"github.com/gocql/gocql"
)

var Session *gocql.Session

func ConnectToCassandra() {
	var err error
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "shwitter"
	Session, err = cluster.CreateSession()
	if err != nil {
		panic(err)
	}
}
