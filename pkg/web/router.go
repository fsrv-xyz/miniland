package web

import (
	"net/http"

	"github.com/gorilla/mux"

	"ref.ci/fsrvcorp/miniland/userland/pkg/web/frontend"
)

func setupRouter() *mux.Router {
	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/addresses", apiAddressesHandler)
	apiRouter.HandleFunc("/files", apiFilesHandler)
	apiRouter.HandleFunc("/processes", apiProcessesHandler)

	frontendRouter := router.PathPrefix("/frontend").Subrouter()
	frontendRouter.HandleFunc("/sse/usage", UsageSSEHandlerBuilder())

	router.PathPrefix("/").Handler(http.FileServer(http.FS(frontend.DistFileSystem())))
	return router
}

func Start(address string) {
	server := http.Server{
		Addr:    address,
		Handler: setupRouter(),
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
