package commands

import (
	"fmt"
	"github.com/Setti7/shwitter/internal/config"
	"github.com/Setti7/shwitter/internal/service"
	"github.com/urfave/cli/v2"
)

var ResetCommand = cli.Command{
	Name:   "reset",
	Usage:  "Resets the database, removing all of its contents.",
	Action: resetAction,
	Flags: []cli.Flag{
		&config.CassandraHostsFlag,
		&config.CassandraKeyspaceFlag,
	}}

func resetAction(ctx *cli.Context) error {
	c := config.NewConfig(ctx)
	service.SetConfig(c)

	sess, err := createCassandraSystemClient(c.Cassandra())
	if err != nil {
		return err
	}
	defer sess.Close()

	// Drop the keyspace and the create it again, running all migrations
	err = dropKeyspace(sess, c.Cassandra())
	if err != nil {
		return err
	}
	l.Infoln(fmt.Sprintf("%s keyspace was dropped.", c.Cassandra().Keyspace))

	err = createKeyspace(sess, c.Cassandra())
	if err != nil {
		return err
	}
	l.Infoln(fmt.Sprintf("Recreated %s keyspace.", c.Cassandra().Keyspace))

	// Running all migrations
	err = runMigrations(c.Cassandra(), 0)
	if err != nil {
		return err
	}

	l.Infoln("All migrations ran successfully.")
	l.Infoln("Reset process is completed.")
	return nil
}
