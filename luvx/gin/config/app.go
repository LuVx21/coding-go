package config

import "github.com/luvx21/coding-go/infra/logs"

type Config struct {
	Server  ServerConfig
	Log     logs.LogConfig
	MySQL   MySQL
	Redis   Redis
	MongoDB MongoDB
	Turso   Turso
}

type ServerConfig struct {
	Port  string
	Debug bool
}

type MySQL struct {
	Host     string
	Port     int
	Username string
	Password string
	Dbname   string
}

type Redis struct {
	Host      string
	Username  string
	Password  string
	Timeout   int
	MaxActive int
	MaxIdle   int
}

type MongoDB struct {
	Uri      string
	Database string
}

type Turso struct {
	Dbname string
	Token  string
}
