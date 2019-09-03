package database

import (
	"os"
)

func setEnvs(h, o, s, d, u, p string) func() {
	if h != "" {
		os.Setenv(ENVHost, h)
	}
	if o != "" {
		os.Setenv(ENVPort, o)
	}
	if s != "" {
		os.Setenv(ENVSSL, s)
	}
	if d != "" {
		os.Setenv(ENVDatabase, d)
	}
	if u != "" {
		os.Setenv(ENVUser, u)
	}
	if p != "" {
		os.Setenv(ENVPassword, p)
	}

	return unsetEnvs
}

func unsetEnvs() {
	os.Unsetenv(ENVHost)
	os.Unsetenv(ENVPort)
	os.Unsetenv(ENVSSL)
	os.Unsetenv(ENVDatabase)
	os.Unsetenv(ENVUser)
	os.Unsetenv(ENVPassword)
}
