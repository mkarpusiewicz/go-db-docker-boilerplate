package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/mkarpusiewicz/go-db-docker-boilerplate/internal/app"
)

var version = "local_build"

func main() {
	app.Startup(version)

	ok := initializationAndHealthchecks()

	if !ok {
		log.Fatal().Msg("exiting...")

		os.Exit(1)
	}

	log.Info().Msg("startup successful")

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := startHTTP(); err != nil {
			log.Fatal().Err(err)
		} else {
			log.Info().Msg("shutting down http server...")
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	log.Info().Str("signal", sig.String()).Msg("notify")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatal().Err(err)
	}
	wg.Wait()

	select {
	case <-ctx.Done():
		log.Info().Msg("exiting with timeout...")
	default:
		log.Info().Msg("exiting...")
	}
}
