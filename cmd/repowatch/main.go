package main

import (
	"os"

	"github.com/moov-io/base/log"
	"github.com/sethlivingston/repowatch"
	"github.com/sethlivingston/repowatch/internal/server"
)

func main() {
	logger := log.NewDefaultLogger().
		Set("app", log.String("repowatch")).
		Set("version", log.String(repowatch.Version))

	env := &server.Environment{
		Logger: logger,
	}

	_, err := server.NewEnvironment(env)
	if err != nil {
		logger.Fatal().LogErrorf("creating environment: %v", err)
		os.Exit(1)
	}

	termListener := server.NewTerminationListener()
	stopServer := env.StartServer(termListener)
	defer stopServer()

	server.AwaitTermination(env.Logger, termListener)
}
