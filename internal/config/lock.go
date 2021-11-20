package config

type LockConfig struct {
	Host string
	Port string
}

var LockDefault = LockConfig{
	Host: "127.0.0.1",
	Port: "6379",
}
