package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

type Options struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseUrl       string `env:"BASE_URL"`
}

var Opt = Options{}

func SetOptions() {
	err := env.Parse(&Opt)
	if err != nil {
		log.Fatal(err)
	}

	if Opt.ServerAddress == "" {
		flag.StringVar(&Opt.ServerAddress, "a", "localhost:8080", "server address")
	}
	if Opt.BaseUrl == "" {
		flag.StringVar(&Opt.BaseUrl, "b", "http://localhost:8080", "base url")
	}

	if Opt.ServerAddress == "" || Opt.BaseUrl == "" {
		flag.Parse()
	}
}
