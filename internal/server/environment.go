package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/moov-io/base/log"
	"github.com/sethlivingston/reponotifications/internal/github"
)

type Environment struct {
	Logger log.Logger
	Config *Config

	GitHubService    github.GitHubService
	GitHubController github.GitHubController

	Router   *mux.Router
	Shutdown func()
}

func NewEnvironment(env *Environment) (*Environment, error) {
	if env == nil {
		env = &Environment{}
	}

	if env.Logger == nil {
		env.Logger = log.NewDefaultLogger()
	}

	if env.Shutdown == nil {
		env.Shutdown = func() {}
	}

	if env.Config == nil {
		config, err := LoadConfig(env.Logger)
		if err != nil {
			return nil, fmt.Errorf("loading config: %v", err)
		}
		env.Config = config
	}

	if env.Router == nil {
		env.Router = mux.NewRouter()
	}

	// Initialize listeners

	if env.Config.GitHub != nil {
		err := env.listenToGitHub()
		if err != nil {
			return nil, fmt.Errorf("setting up github listeners: %v", err)
		}
	}

	// Initialize broadcasters

	if env.Config.Slack != nil {
		err := env.broadcastToSlack()
		if err != nil {
			return nil, fmt.Errorf("setting up slack broadcasters: %v", err)
		}
	}

	return env, nil
}

func (env Environment) StartServer(terminationListener chan error) func() {
	httpServer := &http.Server{
		Addr:         env.Config.Server.BindAddress,
		Handler:      env.Router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		env.Logger.Info().Logf("http server listening on %s", httpServer.Addr)
		err := httpServer.ListenAndServe()
		if err != nil {
			terminationListener <- env.Logger.Fatal().LogError(err).Err()
		}
	}()

	stopServer := func() {
		err := httpServer.Shutdown(context.Background())
		if err != nil {
			env.Logger.Fatal().LogErrorf("unable to shut down properly: %v", err)
		}
	}

	return stopServer
}

func (env Environment) listenToGitHub() error {
	if len(env.Config.GitHub.SigningSecret) == 0 {
		return fmt.Errorf("github signing secret is required")
	}

	if env.GitHubService == nil {
		ghs, err := github.NewGitHubService(env.Logger)
		if err != nil {
			return fmt.Errorf("creating github service: %v", err)
		}
		env.GitHubService = &ghs
	}

	if env.GitHubController == nil {
		ghc := github.NewGitHubController(env.Logger, env.GitHubService)
		env.GitHubController = ghc
	}
	env.GitHubController.AppendRoutes(env.Router)

	return nil
}

func (env Environment) broadcastToSlack() error {
	if len(env.Config.Slack.BotToken) == 0 {
		return fmt.Errorf("slack bot token is required")
	}
	if len(env.Config.Slack.SigningSecret) == 0 {
		return fmt.Errorf("slack signing secret is required")
	}
	return nil
}
