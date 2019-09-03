package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitConfig_Empty(t *testing.T) {
	unsetEnvs()

	_, err := initConfig()

	assert.Error(t, err, "If no environment variables are set there should be an error")
}

func TestInitConfigOK(t *testing.T) {
	h, o, s, d, u, p := "h", "o", "s", "d", "u", "p"
	cancel := setEnvs(h, o, s, d, u, p)
	defer cancel()

	cfg, err := initConfig()

	assert.NoError(t, err, "Environment variables are set there should not be an error")

	assert.Equal(t, h, cfg.host)
	assert.Equal(t, o, cfg.port)
	assert.Equal(t, s, cfg.ssl)
	assert.Equal(t, d, cfg.database)
	assert.Equal(t, u, cfg.user)
	assert.Equal(t, p, cfg.password)
}

func TestInitConfig_NoCredentials(t *testing.T) {
	h, o, s, d := "h", "o", "s", "d"
	cancel := setEnvs(h, o, s, d, "", "")
	defer cancel()

	_, err := initConfig()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "credentials have not been set")
}

func TestInitConfig_NoSettings(t *testing.T) {
	s, d, u, p := "s", "d", "u", "p"
	cancel := setEnvs("", "", s, d, u, p)
	defer cancel()

	_, err := initConfig()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "variables have been set")
}
