// Source:
// https://getstream.io/blog/building-a-performant-api-using-go-and-cassandra/
package main

import (
	"github.com/Setti7/shwitter/Cassandra"
	"github.com/Setti7/shwitter/Shweets"
	"github.com/Setti7/shwitter/Users"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	Cassandra.ConnectToCassandra()
	session := Cassandra.Session
	defer session.Close()

	r := gin.Default()
	r.GET("/healthz", heartbeat)

	r.POST("/shweets/", Shweets.CreateShweet)
	r.GET("/shweets/", Shweets.ListShweets)
	r.GET("/shweets/:uuid", Shweets.GetShweet)

	r.POST("/users/", Users.CreateUser)
	r.GET("/users/", Users.ListUsers)
	r.GET("/users/:uuid", Users.GetUser)

	log.Fatal(r.Run())
}

func heartbeat(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
