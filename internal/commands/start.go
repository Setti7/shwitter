package commands

import (
	"github.com/Setti7/shwitter/internal/api"
	"github.com/Setti7/shwitter/internal/config"
	"github.com/Setti7/shwitter/internal/middleware"
	"github.com/Setti7/shwitter/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
)

var StartCommand = cli.Command{
	Name:    "start",
	Aliases: []string{"up"},
	Usage:   "Starts the web server",
	Action:  startAction,
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

func startAction(ctx *cli.Context) error {
	c := config.NewConfig(ctx)
	service.SetConfig(c)

	service.Init()
	defer service.CleanUp()

	r := gin.Default()
	r.Use(middleware.SessionMiddleware())

	r.GET("/healthz", heartbeat)

	r.POST("/shweets", api.CreateShweet)
	r.GET("/shweets", api.ListShweets)
	r.GET("/shweets/:id", api.GetShweet)

	// TODO: add pagination to ListUsers, ListShweets, ListFollowers and ListFriends
	// TODO: add timeline and userline
	// TODO: add tests
	r.GET("/users", api.ListUsers)
	r.GET("/users/:id", api.GetUser)
	r.POST("/users", api.CreateUser)
	r.GET("/users/me", api.GetCurrentUser)
	r.POST("/users/:id/follow", api.FollowUser)
	r.POST("/users/:id/unfollow", api.UnFollowUser)
	r.GET("/users/:id/followers", api.ListFriendsOrFollowers(false))
	r.GET("/users/:id/friends", api.ListFriendsOrFollowers(true))

	r.POST("/sessions", api.CreateSession)
	r.DELETE("/sessions/:id", api.DeleteSession)
	r.GET("/sessions", api.ListUserSessions)

	// add mentions, then add mentions notifications
	// add chat, after notifications
	// add api rate limiter = 60/min guest 100/min logged

	log.Fatal(r.Run())

	return nil
}

func heartbeat(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
