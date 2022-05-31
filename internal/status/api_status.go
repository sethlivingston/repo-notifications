package status

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/moov-io/base/log"
	"github.com/sethlivingston/repowatch/internal/config"
)

type StatusController interface {
	AppendRoutes(router *mux.Router) *mux.Router
}

func NewStatusController(logger log.Logger, config *config.Config) StatusController {
	return statusController{
		logger: logger,
		config: config,
	}
}

type statusController struct {
	logger log.Logger
	config *config.Config
}

func (c statusController) AppendRoutes(router *mux.Router) *mux.Router {
	router.
		Name("Status.get").
		Methods("GET").
		Path("/status").
		HandlerFunc(c.getStatus)

	return router
}

func (c statusController) getStatus(w http.ResponseWriter, r *http.Request) {
	const Configured string = "Configured\n"
	const NotConfigured string = "Not configured\n"

	s := "REPOWATCH STATUS\n"
	s += "----------------\n"

	s += "\nListeners:\n"

	s += "  GitHub: "
	if c.config.GitHub != nil {
		s += Configured
	} else {
		s += NotConfigured
	}

	s += "\nBroadcasters:\n"

	s += "  Slack: "
	if c.config.Slack != nil {
		s += Configured
	} else {
		s += NotConfigured
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(s))
}
