// +build debug

package integration

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestMain(t *testing.T) {
	e := httpexpect.New(t, serverURL)

	e.GET("/").
		Expect().
		Status(http.StatusOK).
		Text().
		Match("Hello world!")
}

func TestHealthcheck(t *testing.T) {
	e := httpexpect.New(t, serverURL)

	e.GET("/healthcheck").
		Expect().
		Status(http.StatusOK).
		Text().
		Match("Health check OK")
}
