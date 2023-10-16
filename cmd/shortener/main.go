package main

import (
	"fmt"
	"github.com/SergeyGushan/lrn_go_url/cmd/config"
	"github.com/SergeyGushan/lrn_go_url/internal/urlhandlers"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func URLRouter() chi.Router {
	r := chi.NewRouter()
	r.Post("/", urlhandlers.Save)
	r.Get("/{shortCode}", urlhandlers.Get)

	return r
}

func main() {
	config.SetOptions()
	address := fmt.Sprintf("%s:%s", config.Opt.ServerAddress, config.Opt.ServerPort)
	println(address)
	err := http.ListenAndServe(address, URLRouter())
	if err != nil {
		return
	}
}
