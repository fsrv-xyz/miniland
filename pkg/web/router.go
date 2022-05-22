package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

func setupRouter() *mux.Router {
	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/addresses", apiAddressesHandler)
	apiRouter.HandleFunc("/files", apiFilesHandler)

	return router
}

func Start() {
	server := http.Server{
		Addr:    ":8080",
		Handler: setupRouter(),
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
