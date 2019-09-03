package app

import (
	"os"
)

var appENVName = "APP_ENV"

type Environment string

var (
	ProductionENV Environment = "production"
	LocalENV      Environment = "local"
	DebugENV      Environment = "debug"
)

type config struct {
	Env Environment
}

var Config *config

func Init() {
	app, ok := os.LookupEnv(appENVName)

	env := ProductionENV
	if ok {
		env = Environment(app)
	}

	Config = &config{
		Env: env,
	}
}
