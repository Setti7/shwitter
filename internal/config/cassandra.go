package config

import (
	"github.com/urfave/cli/v2"
	"os"
	"strings"
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
	c := &CassandraDefault

	getCassandraConfigFromEnv(c)
	getCassandraConfigFromCLI(c, ctx)

	return c
}

func getCassandraConfigFromEnv(c *CassandraConfig) {
	hosts := os.Getenv(CASSANDRA_HOSTS_ENV)
	if hosts != "" {
		c.Hosts = strings.Split(hosts, ",")
	}

	keyspace := os.Getenv(CASSANDRA_KEYSPACE_ENV)
	if keyspace != "" {
		c.Keyspace = keyspace
	}
}

func getCassandraConfigFromCLI(c *CassandraConfig, ctx *cli.Context) {
	hosts := ctx.StringSlice(CASSANDRA_HOSTS_FLAG_NAME)
	if len(hosts) > 0 {
		c.Hosts = hosts
	}

	keyspace := ctx.String(CASSANDRA_KEYSPACE_FLAG_NAME)
	if keyspace != "" {
		c.Keyspace = keyspace
	}
}
