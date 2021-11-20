package config

type CassandraConfig struct {
	Hosts    []string
	Keyspace string
}

var CassandraDefault = CassandraConfig{
	Hosts:    []string{"127.0.0.1"},
	Keyspace: "shwitter",
}
