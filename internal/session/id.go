package session

import (
	"crypto/rand"
	"fmt"
	"log"
)

// TODO add cassandra support to github.com/gin-contrib/sessions?
func NewID() string {
	b := make([]byte, 24)

	if _, err := rand.Read(b); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", b)
}
