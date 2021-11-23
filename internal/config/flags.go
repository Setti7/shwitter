package config

import (
	"github.com/urfave/cli/v2"
)

var Flags = []cli.Flag{
	// Cassandra settings
	&CassandraHostsFlag,
	&CassandraKeyspaceFlag,

	// Redis settings
	&RedisHostFlag,
}

var CASSANDRA_HOSTS_FLAG_NAME = "cassandra-hosts"
var CassandraHostsFlag = cli.StringSliceFlag{
	Name:     CASSANDRA_HOSTS_FLAG_NAME,
	Required: false,
	Usage:    "A list of the Cassandra hosts.",
}

var CASSANDRA_KEYSPACE_FLAG_NAME = "cassandra-keyspace"
var CassandraKeyspaceFlag = cli.StringFlag{
	Name:     CASSANDRA_KEYSPACE_FLAG_NAME,
	Required: false,
	Usage:    "The Cassandra keyspace for this app.",
}

var REDIS_HOST_FLAG_NAME = "redis-host"
var RedisHostFlag = cli.StringFlag{
	Name:     REDIS_HOST_FLAG_NAME,
	Required: false,
	Usage:    "The host to connect to Redis.",
}
