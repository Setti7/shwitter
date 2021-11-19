// Source:
// https://getstream.io/blog/building-a-performant-api-using-go-and-cassandra/
package main

import (
	"github.com/Setti7/shwitter/Auth"
	"github.com/Setti7/shwitter/api"
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
	r.GET("/healthz", heartbeat)

	r.POST("/shweets/", api.CreateShweet)
	r.GET("/shweets/", api.ListShweets)
	r.GET("/shweets/:id", api.GetShweet)

	r.GET("/users/", api.ListUsers)
	r.GET("/users/:uuid", api.GetUser)

	r.POST("/auth/signup", Auth.SignUp)
	r.POST("/auth/signin", Auth.SignIn)

	log.Fatal(r.Run())
}

func heartbeat(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
