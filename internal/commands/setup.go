package commands

import (
	"fmt"
	"github.com/Setti7/shwitter/internal/config"
	"github.com/Setti7/shwitter/internal/log"
	"github.com/Setti7/shwitter/internal/service"
	"github.com/gocql/gocql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/cassandra"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/urfave/cli/v2"
)

var l = log.Log

var SetupCommand = cli.Command{
	Name:   "setup",
	Usage:  "Setups the dependencies for the web server to start",
	Action: setupAction,
	Flags: []cli.Flag{
		// Cassandra settings
		&cli.StringSliceFlag{
			Name:  "cassandra-hosts",
			Value: cli.NewStringSlice(config.CassandraDefault.Hosts...),
			Usage: "A list of the Cassandra hosts.",
		},
		&cli.StringFlag{
			Name:  "cassandra-keyspace",
			Value: config.CassandraDefault.Keyspace,
			Usage: "The Cassandra keyspace for this app.",
		},
		// Redis settings
		&cli.StringFlag{
			Name:  "redis-host",
			Value: config.LockDefault.Host,
			Usage: "The host to connect to Redis.",
		},
	},
}

func setupAction(ctx *cli.Context) error {
	c := config.NewConfig(ctx)
	service.SetConfig(c)

	// Create a new cassandra client to the system keyspace so we can create our own keyspace
	cluster := gocql.NewCluster(c.Cassandra().Hosts...)
	cluster.Keyspace = "system"

	sess, err := cluster.CreateSession()
	if err != nil {
		log.LogError("setupAction", "Could not connect to the system keyspace on Cassandra", err)
		return err
	}

	// Create the required keyspace
	err = sess.Query(fmt.Sprintf("CREATE KEYSPACE %s WITH replication = "+
		"{'class': 'SimpleStrategy', 'replication_factor': 1};", c.Cassandra().Keyspace)).Exec()
	if err != nil {
		log.LogError("setupAction", fmt.Sprintf("Could not create the %s keyspace.",
			c.Cassandra().Keyspace), err)
		return err
	}

	sess.Close()
	l.Infoln(fmt.Sprintf("Keyspace %s was created successfully.", c.Cassandra().Keyspace))

	// Running all migrations
	d, err := cassandra.WithInstance(service.Cassandra(), &cassandra.Config{KeyspaceName: c.Cassandra().Keyspace,
		MultiStatementEnabled: true})
	if err != nil {
		log.LogError("setupAction", fmt.Sprintf("Could not connect to the %s keyspace.",
			c.Cassandra().Keyspace), err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", c.Cassandra().Keyspace, d)
	if err != nil {
		log.LogError("setupAction", "Could not find the migrations folder.", err)
		return err
	}

	err = m.Up()
	if err != nil {
		log.LogError("setupAction", "Could not run migrations.", err)
		return err
	}

	l.Infoln("All migration ran successfully.")
	l.Infoln("Setup process is completed.")
	return nil
}
