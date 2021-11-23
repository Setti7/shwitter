package commands

import (
	"fmt"
	"github.com/Setti7/shwitter/internal/config"
	"github.com/Setti7/shwitter/internal/service"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/urfave/cli/v2"
	"strconv"
)

var MigrateCommand = cli.Command{
	Name:  "migrate",
	Usage: "Migrates the database.",
	Description: "The argument is the number of migrations that will be run. If the argument is left empty, all " +
		"migrations will be applied. If it's negative, then the last migrations will be rolled back.",
	Action:    migrateAction,
	ArgsUsage: "[number of migrations]",
	Flags: []cli.Flag{
		&config.CassandraHostsFlag,
		&config.CassandraKeyspaceFlag,
	},
}

func migrateAction(ctx *cli.Context) (err error) {
	numOfMigrations := 0 // by default, will run all migrations

	// Parse the argument if there is any
	if ctx.NArg() == 1 {
		arg := ctx.Args().Get(0)
		numOfMigrations, err = strconv.Atoi(arg)
		if err != nil {
			l.Error("The number of migrations needs to be a valid (positive or negative) integer.")
			return nil
		}
	}

	c := config.NewConfig(ctx)
	service.SetConfig(c)

	err = createKeyspace(c.Cassandra())
	if err != nil {
		return err
	}
	l.Infoln(fmt.Sprintf("Using %s keyspace.", c.Cassandra().Keyspace))

	err = runMigrations(c.Cassandra(), numOfMigrations)
	if err != nil {
		return err
	}

	if numOfMigrations == 0 {
		l.Infoln("All migrations ran successfully.")
	} else {
		l.Infoln(fmt.Sprintf("%d migration(s) ran successfully.", numOfMigrations))
	}

	l.Infoln("Success.")
	return nil
}
