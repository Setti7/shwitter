package commands

import (
	"log"
	"net/http"
	"time"

	"github.com/Setti7/shwitter/internal/config"
	"github.com/Setti7/shwitter/internal/follow"
	"github.com/Setti7/shwitter/internal/middleware"
	"github.com/Setti7/shwitter/internal/service"
	"github.com/Setti7/shwitter/internal/session"
	"github.com/Setti7/shwitter/internal/shweets"
	"github.com/Setti7/shwitter/internal/users"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
)

var StartCommand = cli.Command{
	Name:    "start",
	Aliases: []string{"up"},
	Usage:   "Starts the web server",
	Action:  startAction,
	Flags:   config.Flags,
}

func startAction(ctx *cli.Context) error {
	c := config.NewConfig(ctx)
	service.SetConfig(c)

	service.Init()
	defer service.CleanUp()

	usersRepo := users.NewCassandraRepository(service.Cassandra())
	usersService := users.NewService(usersRepo, service.Lock())

	sessRepo := session.NewCassandraRepository(service.Cassandra())
	sessService := session.NewService(sessRepo, usersRepo)

	followRepo := follow.NewCassandraRepository(service.Cassandra(), usersRepo)
	followService := follow.NewService(followRepo)

	shweetRepo := shweets.NewCassandraRepository(service.Cassandra(), usersRepo)
	shweetService  := shweets.NewService(shweetRepo)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Content-Length", "Content-Type", "accept", "origin", "Cache-Control", "X-Session-Token"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(middleware.SessionMiddleware(sessService))
	r.Use(middleware.UserMiddleware(usersService))

	r.GET("/healthz", heartbeat)

	r.POST("/shweets", api.CreateShweet)
	r.GET("/shweets", api.ListShweets)
	r.GET("/shweets/:id", api.GetShweet)
	r.POST("/shweets/:id/like", api.LikeOrUnlikeShweet)

	// TODO: add tests, interface and channels
	r.GET("/users", api.ListUsers)
	r.GET("/users/:id", api.GetUser)
	r.GET("/users/:id/profile", api.GetUserProfile)
	r.POST("/users", api.CreateUser)
	r.GET("/users/me", api.GetCurrentUser)

	r.GET("/users/:id/follow", api.IsFollowingUser)
	r.POST("/users/:id/follow", api.FollowOrUnfollowUser)
	r.GET("/users/:id/followers", api.ListFriendsOrFollowers(false))
	r.GET("/users/:id/friends", api.ListFriendsOrFollowers(true))

	r.POST("/sessions", api.CreateSession)
	r.DELETE("/sessions/:id", api.DeleteSession)
	r.GET("/sessions", api.ListUserSessions)

	r.GET("/timeline", api.GetTimelineForCurrentUser)
	r.GET("/userline/:id", api.GetUserLine)

	// add mentions, then add mentions notifications
	// add chat, after notifications
	// add api rate limiter = 60/min guest 100/min logged

	log.Fatal(r.Run())

	return nil
}

func heartbeat(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
