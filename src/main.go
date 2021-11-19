// Source:
// https://getstream.io/blog/building-a-performant-api-using-go-and-cassandra/
package main

import (
	"github.com/Setti7/shwitter/api"
	"github.com/Setti7/shwitter/middleware"
	"github.com/Setti7/shwitter/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// TODO: try https://github.com/scylladb/gocqlx
// TODO: add tests

func main() {
	service.Init()
	defer service.CleanUp()

	r := gin.Default()
	r.Use(middleware.CurrentUserMiddleware())

	r.GET("/healthz", heartbeat)

	r.POST("/shweets/", api.CreateShweet)
	r.GET("/shweets/", api.ListShweets)
	r.GET("/shweets/:id", api.GetShweet)

	// TODO acl
	r.GET("/users/", api.ListUsers)
	r.GET("/users/:uuid", api.GetUser)
	r.POST("/users/", api.CreateUser)
	r.GET("/users/me", api.GetCurrentUser)

	r.POST("/sessions/", api.CreateSession)
	// TODO delete session

	log.Fatal(r.Run())
}

func heartbeat(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
