package github

import "github.com/moov-io/base/log"

type GitHubService interface{}

func NewGitHubService(logger log.Logger) (GitHubService, error) {
	return &gitHubService{
		logger: logger,
	}, nil
}

type gitHubService struct {
	logger log.Logger
}
