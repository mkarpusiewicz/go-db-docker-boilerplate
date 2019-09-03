package app

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger

func Startup(version string) {
	envSource := "file"
	if err := godotenv.Load(); err != nil {
		envSource = "vars"
	}

	Init()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if Config.Env != ProductionENV {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"})
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	Logger = log.Logger

	log.Info().
		Str("version", version).
		Str("env_vars", envSource).
		Str("app_env", string(Config.Env)).
		Msg("startup initialization")
}
