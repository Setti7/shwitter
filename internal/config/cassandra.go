package config

import (
	"github.com/urfave/cli/v2"
)

type CassandraConfig struct {
	Hosts    []string
	Keyspace string
}

var CassandraDefault = CassandraConfig{
	Hosts:    []string{"127.0.0.1"},
	Keyspace: "shwitter",
}

func (c *Config) Cassandra() *CassandraConfig {
	return c.cassandra
}

func NewCassandraConfig(ctx *cli.Context) *CassandraConfig {
	c := &CassandraConfig{}

	if ctx == nil {
		return &CassandraDefault
	}

	c.Hosts = ctx.StringSlice("cassandra-hosts")
	c.Keyspace = ctx.String("cassandra-keyspace")

	return c
}
