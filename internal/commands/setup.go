package commands

import (
	"fmt"
	"github.com/Setti7/shwitter/internal/config"
	"github.com/Setti7/shwitter/internal/log"
	"github.com/Setti7/shwitter/internal/service"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/cassandra"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/urfave/cli/v2"
)

var SetupCommand = cli.Command{
	Name:   "setup",
	Usage:  "Setups the dependencies for the web server to start",
	Action: setupAction,
	Flags:  config.Flags,
}

func setupAction(ctx *cli.Context) error {
	c := config.NewConfig(ctx)
	service.SetConfig(c)

	err := createKeyspace(c.Cassandra())
	if err != nil {
		return err
	}
	l.Infoln(fmt.Sprintf("Using %s keyspace.", c.Cassandra().Keyspace))

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

	l.Infoln("All migrations ran successfully.")
	l.Infoln("Setup process is completed.")
	return nil
}
