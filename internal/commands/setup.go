package commands

import (
	"fmt"
	"github.com/Setti7/shwitter/internal/config"
	"github.com/Setti7/shwitter/internal/service"
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
	err = runMigrations(c.Cassandra(), 0)
	if err != nil {
		return err
	}

	l.Infoln("All migrations ran successfully.")
	l.Infoln("Setup process is completed.")
	return nil
}
