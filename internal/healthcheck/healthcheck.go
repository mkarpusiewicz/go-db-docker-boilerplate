package healthcheck

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var Timeout = 1 * time.Second

type HealthcheckReply struct {
	Success bool     `json:"success"`
	TimeMs  int32    `json:"time_ms"`
	Errors  []string `json:"errors,omitempty"`
}

func Handler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	start := time.Now()
	errors, ok := Perform(ctx)
	elapsed := time.Since(start)
	ms := int32(elapsed / time.Millisecond)

	errMsgs := make([]string, len(errors))
	for i, err := range errors {
		errMsgs[i] = err.Error()
	}

	response := &HealthcheckReply{
		Success: ok,
		TimeMs:  ms,
		Errors:  errMsgs,
	}

	status := http.StatusOK
	if !ok {
		status = http.StatusServiceUnavailable
	}

	c.JSON(status, response)
}
