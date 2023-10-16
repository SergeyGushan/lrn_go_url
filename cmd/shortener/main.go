package main

import (
	"github.com/SergeyGushan/lrn_go_url/cmd/config"
	"github.com/SergeyGushan/lrn_go_url/internal/urlhandlers"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	config.SetOptions()
	err := http.ListenAndServe(config.Opt.ServerAddress, URLRouter())
	if err != nil {
		return
	}
}

func URLRouter() chi.Router {
	r := chi.NewRouter()
	r.Post("/", urlhandlers.Save)
	r.Get("/{shortCode}", urlhandlers.Get)

	return r
}
