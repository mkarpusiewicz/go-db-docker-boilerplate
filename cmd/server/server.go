package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	
	"github.com/mkarpusiewicz/go-db-docker-boilerplate/internal/app"
	"github.com/mkarpusiewicz/go-db-docker-boilerplate/internal/healthcheck"
)

var s *http.Server

func startHTTP() error {
	if app.Config.Env == app.ProductionENV {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	if app.Config.Env == app.ProductionENV {
		r.Use(logger.SetLogger())
	} else {
		r.Use(gin.Logger())
	}

	r.GET("/", func(c *gin.Context) { c.String(http.StatusOK, "Hello world!") })
	r.GET("/healthcheck", healthcheck.Handler)

	if app.Config.Env == app.LocalENV {
		d := r.Group("/debug")
		d.GET("/sleep", func(c *gin.Context) {
			time.Sleep(5 * time.Second)
			c.String(http.StatusOK, "Sleep test!")
		})
	}

	p := os.Getenv("SERVER_PORT")
	port, err := strconv.Atoi(p)
	if err != nil {
		log.Fatal().Str("SERVER_PORT", p).Msg("Invalid port value")
	}
	log.Info().Int("port", port).Msg("http server starting")

	s = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
