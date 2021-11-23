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
	"strconv"
)

var MigrateCommand = cli.Command{
	Name:   "migrate",
	Usage:  "Migrates the database",
	Action: migrateAction,
	//ArgsUsage: TODO
	Flags: []cli.Flag{
		&config.CassandraHostsFlag,
		&config.CassandraKeyspaceFlag,
	},
}

func migrateAction(ctx *cli.Context) (err error) {
	var numOfMigrations int
	runAllMigrations := ctx.NArg() == 0

	// Parse the argument if there is any
	if !runAllMigrations {
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

	// Running migrations
	d, err := cassandra.WithInstance(service.Cassandra(), &cassandra.Config{KeyspaceName: c.Cassandra().Keyspace,
		MultiStatementEnabled: true})
	if err != nil {
		log.LogError("migrateAction", fmt.Sprintf("Could not connect to the %s keyspace.",
			c.Cassandra().Keyspace), err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", c.Cassandra().Keyspace, d)
	if err != nil {
		log.LogError("migrateAction", "Could not find the migrations folder.", err)
		return err
	}

	if runAllMigrations {
		err = m.Up()
	} else {
		err = m.Steps(numOfMigrations)
	}

	if err != nil {
		log.LogError("migrateAction", "Could not run migrations.", err)
		return err
	}

	if runAllMigrations {
		l.Infoln("All migrations ran successfully.")
	} else {
		l.Infoln(fmt.Sprintf("%d migration(s) ran successfully.", numOfMigrations))
	}

	l.Infoln("Success.")
	return nil
}
