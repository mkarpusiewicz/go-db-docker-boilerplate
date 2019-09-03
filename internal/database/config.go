package database

import (
	"fmt"
	"os"
)

//Environment variable names for database config
const (
	ENVHost     = "POSTGRES_HOST"
	ENVPort     = "POSTGRES_PORT"
	ENVSSL      = "POSTGRES_SSL"
	ENVDatabase = "POSTGRES_DB"
	ENVUser     = "POSTGRES_USER"
	ENVPassword = "POSTGRES_PASSWORD"
)

type databaseConfig struct {
	host     string
	port     string
	ssl      string
	database string
	user     string
	password string
}

func initConfig() (databaseConfig, error) {
	cfg := databaseConfig{
		host:     os.Getenv(ENVHost),
		port:     os.Getenv(ENVPort),
		ssl:      os.Getenv(ENVSSL),
		database: os.Getenv(ENVDatabase),
		user:     os.Getenv(ENVUser),
		password: os.Getenv(ENVPassword),
	}

	if cfg.user == "" || cfg.password == "" {
		return cfg, fmt.Errorf("postgres credentials have not been set")
	}

	if cfg.host == "" || cfg.port == "" || cfg.database == "" {
		return cfg, fmt.Errorf("not all postgres database server environment variables have been set")
	}

	return cfg, nil
}
