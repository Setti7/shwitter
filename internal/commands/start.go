package commands

import (
	"log"
	"net/http"
	"time"

	"github.com/Setti7/shwitter/internal/api"
	"github.com/Setti7/shwitter/internal/config"
	"github.com/Setti7/shwitter/internal/follow"
	"github.com/Setti7/shwitter/internal/middleware"
	"github.com/Setti7/shwitter/internal/service"
	"github.com/Setti7/shwitter/internal/session"
	"github.com/Setti7/shwitter/internal/shweets"
	"github.com/Setti7/shwitter/internal/timeline"
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

	cass := service.Cassandra()
	lock := service.Lock()

	usersRepo := users.NewCassandraRepository(cass)
	usersService := users.NewService(usersRepo, lock)

	sessRepo := session.NewCassandraRepository(cass)
	sessService := session.NewService(sessRepo, usersService)

	// TODO instead of using userRepo on followRepo, we should use userSvc in followService
	followRepo := follow.NewCassandraRepository(cass, usersRepo)
	followService := follow.NewService(followRepo)

	// TODO instead of using userRepo on shweetRepo, we should use userSvc in shweetService
	shweetRepo := shweets.NewCassandraRepository(cass, usersRepo)
	shweetService := shweets.NewService(shweetRepo)

	// TODO instead of using userRepo and shweetRepo on timelineRepo, we should use userSvc and shweetSvc in timelineService
	timelineRepo := timeline.NewCassandraRepository(cass, usersRepo, shweetRepo)
	timelineService := timeline.NewService(timelineRepo, shweetService)

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

	api.MakeUsersHandlers(r, usersService)
	api.MakeTimelineHandlers(r, timelineService)
	api.MakeShweetsHandlers(r, shweetService)
	api.MakeSessionHandlers(r, sessService)
	api.MakeFollowHandlers(r, followService)

	// TODO: add tests, interface and channels
	// add mentions, then add mentions notifications
	// add chat, after notifications
	// add api rate limiter = 60/min guest 100/min logged

	log.Fatal(r.Run())
	return nil
}

func heartbeat(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
