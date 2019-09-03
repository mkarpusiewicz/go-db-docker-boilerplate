package app

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoEnv_Production(t *testing.T) {
	os.Unsetenv(appENVName)
	Init()

	assert.Equal(t, ProductionENV, Config.Env)
}

func TestDebugEnv_Debug(t *testing.T) {
	os.Setenv(appENVName, string(LocalENV))
	defer os.Unsetenv(appENVName)
	Init()

	assert.Equal(t, LocalENV, Config.Env)
}
