package commands

import (
	"fmt"
	"github.com/Setti7/shwitter/internal/config"
	"github.com/Setti7/shwitter/internal/log"
	"github.com/gocql/gocql"
)

func createKeyspace(c *config.CassandraConfig) (err error) {
	// Create a new cassandra client to the system keyspace, so we can create our own keyspace
	cluster := gocql.NewCluster(c.Hosts...)
	cluster.Keyspace = "system"

	sess, err := cluster.CreateSession()
	if err != nil {
		log.LogError("createKeyspace", "Could not connect to the system keyspace on Cassandra", err)
		return
	}

	// Create the required keyspace
	err = sess.Query(fmt.Sprintf("CREATE KEYSPACE IF NOT EXISTS %s WITH replication = "+
		"{'class': 'SimpleStrategy', 'replication_factor': 1};", c.Keyspace)).Exec()
	if err != nil {
		log.LogError("createKeyspace",
			fmt.Sprintf("Could not create the %s keyspace.", c.Keyspace), err)
		return
	}

	sess.Close()
	return
}
