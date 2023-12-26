package main

import (
	"github.com/SergeyGushan/lrn_go_url/cmd/config"
	"github.com/SergeyGushan/lrn_go_url/internal/app/migrations"
	"github.com/SergeyGushan/lrn_go_url/internal/database"
	"github.com/SergeyGushan/lrn_go_url/internal/handlers"
	"github.com/SergeyGushan/lrn_go_url/internal/logger"
	"github.com/SergeyGushan/lrn_go_url/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func URLRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(logger.Handler)
	r.Post("/", handlers.Save)
	r.Get("/ping", handlers.PingDB)
	r.Post("/api/shorten", handlers.Shorten)
	r.Post("/api/shorten/batch", handlers.BatchCreate)
	r.Get("/{shortCode}", handlers.Get)

	return r
}

func main() {
	err := logger.Initialize("Info")
	if err != nil {
		panic(err)
	}

	config.SetOptions()

	if config.Opt.DatabaseDSN != "" {
		database.Connect()
		migrations.Handle()
		storage.Service = storage.NewDatabaseStorage(database.Client())
	} else if config.Opt.FileStoragePath != "" {
		storage.Service, _ = storage.NewJSONStorage(config.Opt.FileStoragePath)
	} else {
		storage.Service, _ = storage.NewStorage()
	}

	err = http.ListenAndServe(config.Opt.ServerAddress, URLRouter())
	if err != nil {
		panic(err)
	}
}
