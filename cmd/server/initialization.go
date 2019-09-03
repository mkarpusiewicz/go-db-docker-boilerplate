package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kamilsk/breaker"
	"github.com/kamilsk/retry/v4"
	"github.com/kamilsk/retry/v4/backoff"
	"github.com/kamilsk/retry/v4/strategy"
	"github.com/rs/zerolog/log"

	"github.com/mkarpusiewicz/go-db-docker-boilerplate/internal/database"
	"github.com/mkarpusiewicz/go-db-docker-boilerplate/internal/healthcheck"
)

const startupTimeoutSeconds = 10

func initializationAndHealthchecks() bool {
	msgChan := make(chan string)
	errChan := make(chan string)

	timeoutBreakerFactory := func() retry.BreakCloser {
		return breaker.BreakByTimeout(startupTimeoutSeconds * time.Second)
	}
	retryStrategyFactory := func() func(attempt uint, err error) bool {
		return strategy.Backoff(backoff.Exponential(500*time.Millisecond, 2))
	}
	errorLogStrategyFactory := func(name string) func(attempt uint, err error) bool {
		return func(attempt uint, err error) bool {
			if attempt == 0 {
				return true
			}

			if err != nil {
				log.Warn().
					Err(err).
					Uint("attempt", attempt).
					Msgf("%s", name)
			}

			return true
		}
	}

	var pwg sync.WaitGroup

	pwg.Add(1)
	go func() {
		defer pwg.Done()
		if err := retry.Retry(timeoutBreakerFactory(), func(uint) error {
			if err := database.InitDB(); err != nil {
				return fmt.Errorf("cannot initialize db: %v", err)
			}
			return nil
		}, retryStrategyFactory(), errorLogStrategyFactory("database initialization")); err != nil {
			errChan <- "database initialization failed"
		} else {
			msgChan <- "database initialization successful"
		}
	}()

	var cwg sync.WaitGroup

	cwg.Add(1)
	go func() {
		defer cwg.Done()
		for msg := range msgChan {
			log.Info().Msg(msg)
		}
	}()

	var startupErrors []string
	cwg.Add(1)
	go func() {
		defer cwg.Done()
		for err := range errChan {
			startupErrors = append(startupErrors, err)
		}
	}()

	pwg.Wait()
	close(msgChan)
	close(errChan)
	cwg.Wait()

	if len(startupErrors) > 0 {
		log.Fatal().
			Strs("errors", startupErrors).
			Int("timeout", startupTimeoutSeconds).
			Msg("startup failed")

		return false
	}

	// perform healtchecks to be sure all dependencies are ok
	ctx, cancel := context.WithTimeout(context.Background(), startupTimeoutSeconds*time.Second)
	defer cancel()
	healtcheckErrors, ok := healthcheck.Perform(ctx)

	if !ok {
		msgs := make([]string, len(healtcheckErrors))
		for i, err := range healtcheckErrors {
			msgs[i] = err.Error()
		}
		log.
			Fatal().
			Strs("errors", msgs).
			Msg("healthchecks failed")

		return false
	}

	log.Info().Msg("healthchecks successful")

	return true
}
