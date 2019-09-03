package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/mkarpusiewicz/go-db-docker-boilerplate/internal/app"
	"github.com/mkarpusiewicz/go-db-docker-boilerplate/internal/healthcheck"
)

const timeout = 3 * time.Second

var version = "local_build"

func main() {
	app.Startup(version)

	addr := os.Getenv("SERVER_URL") + "/healthcheck"

	logger := app.Logger.
		With().
		Str("server_url", addr).
		Logger()

	startTime := time.Now()
	reply, err := callHealthcheck(addr)
	duration := time.Since(startTime)

	if err != nil {
		logger.Error().
			Err(err).
			Dur("time_taken", duration).
			Msg("healtcheck failed")

		os.Exit(1)
	}

	logger.Info().
		Interface("healthcheck", reply).
		Dur("request_duration_ms", duration).
		Msg("healtcheck ok")

	os.Exit(0)
}

func callHealthcheck(url string) (*healthcheck.HealthcheckReply, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "request creation failed")
	}

	ctx, cancel := context.WithTimeout(req.Context(), timeout)
	defer cancel()

	req = req.WithContext(ctx)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "request sending failed")
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "reading response failed")
	}

	reply := &healthcheck.HealthcheckReply{}
	if err := json.Unmarshal(data, reply); err != nil {
		return nil, errors.Wrap(err, "unmarshalling json failed")
	}

	return reply, nil
}
