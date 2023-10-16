package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

type Options struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

var Opt = Options{}

func SetOptions() {

	flag.StringVar(&Opt.ServerAddress, "a", "localhost:8080", "server address")
	flag.StringVar(&Opt.BaseURL, "b", "http://localhost:8080", "base url")
	flag.Parse()

	err := env.Parse(&Opt)
	if err != nil {
		log.Fatal(err)
	}
}
