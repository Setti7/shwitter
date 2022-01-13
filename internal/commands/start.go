package commands

import (
	"github.com/Setti7/shwitter/internal/api"
	"github.com/Setti7/shwitter/internal/config"
	"github.com/Setti7/shwitter/internal/middleware"
	"github.com/Setti7/shwitter/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
	"time"
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

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Content-Length", "Content-Type", "accept", "origin", "Cache-Control", "X-Session-Token"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(middleware.SessionMiddleware())

	r.GET("/healthz", heartbeat)

	r.POST("/shweets", api.CreateShweet)
	r.GET("/shweets", api.ListShweets)
	r.GET("/shweets/:id", api.GetShweet)

	// TODO: add tests, interface and channels
	r.GET("/users", api.ListUsers)
	r.GET("/users/:id", api.GetUser)
	r.POST("/users", api.CreateUser)
	r.GET("/users/me", api.GetCurrentUser)

	r.POST("/users/:id/follow", api.FollowUser)
	r.GET("/users/:id/follow", api.IsFollowingUser)
	r.POST("/users/:id/unfollow", api.UnFollowUser)
	r.GET("/users/:id/followers", api.ListFriendsOrFollowers(false))
	r.GET("/users/:id/friends", api.ListFriendsOrFollowers(true))

	r.POST("/sessions", api.CreateSession)
	r.DELETE("/sessions/:id", api.DeleteSession)
	r.GET("/sessions", api.ListUserSessions)

	// TODO:
	//  Add a test where 50K users are created, with each of them following other users (following the twitter
	//  followers dataset distribution), and then check if the most famous user can shweet a message for everyone.
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
