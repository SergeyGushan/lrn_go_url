package main

import (
	"github.com/SergeyGushan/lrn_go_url/cmd/config"
	"github.com/SergeyGushan/lrn_go_url/internal/logger"
	"github.com/SergeyGushan/lrn_go_url/internal/storage"
	"github.com/SergeyGushan/lrn_go_url/internal/urlhandlers"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func URLRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(logger.Handler)
	r.Post("/", urlhandlers.Save)
	r.Post("/api/shorten", urlhandlers.Shorten)
	r.Get("/{shortCode}", urlhandlers.Get)

	return r
}

func main() {
	err := logger.Initialize("Info")
	if err != nil {
		panic(err)
	}

	config.SetOptions()

	storage.URLStore, err = storage.NewURL(config.Opt.FileStoragePath)
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(config.Opt.ServerAddress, URLRouter())
	if err != nil {
		panic(err)
	}
}
