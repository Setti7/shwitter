// Source:
// https://getstream.io/blog/building-a-performant-api-using-go-and-cassandra/
package main

import (
	"fmt"
	"github.com/Setti7/shwitter/internal/commands"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
)

var version = "development"

//var log = event.log
// TODO
//   [ ] Add error handling (like done on query.FollowUser)
//	 [ ] Add logging

func main() {
	app := cli.NewApp()
	app.Name = "Shwitter"
	app.HelpName = filepath.Base(os.Args[0])
	app.Usage = "Shitpost like there is no tomorrow"
	app.Description = "Shwitter is like twitter but where you have fun instead of being " +
		"pissed of by other people's stupidity"
	app.Version = version
	app.Copyright = "(c) 2021 André Niero Setti <ansetti7@gmail.com>"
	app.EnableBashCompletion = true

	// TODO add new commands
	//  [ ] Setup initial environment -> create keyspace and migrate
	//  [ ] Start docker-compose
	//  [ ] Start with hot-reload (dev)
	//  [ ] migrate down
	//  [ ] clear database
	//  [ ] backup database
	//  [ ] restore database
	app.Commands = []*cli.Command{
		&commands.StartCommand,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf(err.Error())
	}
}
