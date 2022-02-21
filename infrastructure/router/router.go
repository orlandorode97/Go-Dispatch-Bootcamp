package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/orlandorode97/go-disptach/infrastructure/middleware"
	"github.com/orlandorode97/go-disptach/interface/controller"
)

// NewRouter returns a mux router with the needed resources
func NewRouter(c controller.AppController) *mux.Router {
	router := mux.NewRouter()
	router.Use(middleware.HeadersMiddleware)

	router = router.PathPrefix("/api/v1/").Subrouter()

	// /api/v1/words/
	router.HandleFunc("/words/", c.Words.GetWords).
		Queries("url", "{url}").
		Methods(http.MethodGet)

	// /api/v1/words/{id}/
	router.HandleFunc("/words/{id:[0-9a-zA-Z\\W]+|}/", c.Words.GetWordsFromCSV).
		Methods(http.MethodGet)

	// /api/v1/definitions-csv/
	router.HandleFunc("/words-csv/", c.Words.GetConcurrentWords).Queries(
		"type", "{type}",
		"items", "{items}",
		"items_per_workers", "{items_per_workers}",
	).Methods(http.MethodGet)
	return router
}
