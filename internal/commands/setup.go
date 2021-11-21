package commands

import (
	"github.com/Setti7/shwitter/internal/log"
	"github.com/Setti7/shwitter/internal/service"
	"github.com/urfave/cli/v2"
)

var l = log.Log

var SetupCommand = cli.Command{
	Name:   "setup",
	Usage:  "Setups the dependencies for the web server to start",
	Action: setupAction,
}

func setupAction(_ *cli.Context) error {
	// FIXME needs to load a different configuration for the cassandra module, because we need to use the system
	// keyspace
	err := service.Cassandra().Query("CREATE KEYSPACE IF NOT EXISTS shwitter WITH replication = " +
		"{'class': 'SimpleStrategy', 'replication_factor': 1};").Exec()

	if err != nil {
		log.LogError("setupAction", "An error occurred when setupping the environment.", err)
	} else {
		l.Infoln("Environment was successfully setup.")
	}

	return nil
}
