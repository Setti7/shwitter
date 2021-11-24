package commands

import (
	"fmt"
	"github.com/Setti7/shwitter/internal/config"
	"github.com/Setti7/shwitter/internal/log"
	"github.com/Setti7/shwitter/internal/service"
	"github.com/gocql/gocql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/cassandra"
)

func createCassandraSystemClient(c *config.CassandraConfig) (sess *gocql.Session, err error) {
	cluster := gocql.NewCluster(c.Hosts...)
	cluster.Keyspace = "system"

	sess, err = cluster.CreateSession()
	if err != nil {
		log.LogError("createCassandraSystemClient", "Could not connect to the system keyspace on Cassandra", err)
	}
	return
}

func createKeyspace(sess *gocql.Session, c *config.CassandraConfig) (err error) {
	err = sess.Query(fmt.Sprintf("CREATE KEYSPACE IF NOT EXISTS %s WITH replication = "+
		"{'class': 'SimpleStrategy', 'replication_factor': 1};", c.Keyspace)).Exec()
	if err != nil {
		log.LogError("createKeyspace",
			fmt.Sprintf("Could not create the %s keyspace.", c.Keyspace), err)
	}
	return
}

func dropKeyspace(sess *gocql.Session, c *config.CassandraConfig) (err error) {
	err = sess.Query(fmt.Sprintf("DROP KEYSPACE %s", c.Keyspace)).Exec()
	if err != nil {
		log.LogError("dropKeyspace",
			fmt.Sprintf("Could not drop the %s keyspace.", c.Keyspace), err)
	}
	return
}

func runMigrations(c *config.CassandraConfig, n int) (err error) {
	d, err := cassandra.WithInstance(service.Cassandra(), &cassandra.Config{KeyspaceName: c.Keyspace,
		MultiStatementEnabled: true})
	if err != nil {
		log.LogError("runMigrations", fmt.Sprintf("Could not connect to the %s keyspace.",
			c.Keyspace), err)
		return
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", c.Keyspace, d)
	if err != nil {
		log.LogError("runMigrations", "Could not find the migrations folder.", err)
		return
	}

	// Run migrations
	if n == 0 {
		err = m.Up()
	} else {
		err = m.Steps(n)
	}

	if err != nil {
		log.LogError("runMigrations", "Could not run migrations.", err)
	}
	return
}
