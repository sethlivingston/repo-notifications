package github

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/moov-io/base/log"
)

type GitHubController interface {
	AppendRoutes(router *mux.Router) *mux.Router
}

func NewGitHubController(logger log.Logger, service GitHubService) GitHubController {
	return gitHubController{
		logger:  logger,
		service: service,
	}
}

type gitHubController struct {
	logger  log.Logger
	service GitHubService
}

func (c gitHubController) AppendRoutes(router *mux.Router) *mux.Router {
	router.
		Name("GitHub.webhooks").
		Methods("POST").
		Path("/github/webhooks").
		HandlerFunc(c.handleWebhook)

	return router
}

func (c gitHubController) handleWebhook(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement GitHub webhook handler
	w.WriteHeader(http.StatusNoContent)
}
