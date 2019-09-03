// +build integration

package integration

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/namsral/flag"
)

var (
	serverURL string
	cwd       string
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(".env file not found, using environment variable")
	} else {
		fmt.Println("Using environment variable from .env file")
	}

	flag.StringVar(&serverURL, "server-url", "", "set http rest server address")
	flag.StringVar(&cwd, "cwd", "", "set cwd")
	flag.Parse()

	if serverURL == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Println()

	fmt.Println("Server url:", serverURL)

	if cwd != "" {
		if err := os.Chdir(cwd); err != nil {
			panic(fmt.Sprintf("Chdir error: %v", err))
		}
	}
	fmt.Println("CWD:", cwd)
	fmt.Println()
}
