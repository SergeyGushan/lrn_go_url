package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

type Options struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
	ServerPort    string `env:"SERVER_PORT"`
}

var Opt = Options{}

func SetOptions() {
	err := env.Parse(&Opt)
	if err != nil {
		log.Fatal(err)
	}

	if Opt.ServerAddress == "" {
		flag.StringVar(&Opt.ServerAddress, "a", "localhost", "server address")
	}

	if Opt.ServerPort == "" {
		flag.StringVar(&Opt.ServerPort, "server-port", "8080", "server port")
	}

	if Opt.BaseURL == "" {
		flag.StringVar(&Opt.BaseURL, "b", "http://localhost:8080", "base url")
	}

	if Opt.ServerAddress == "" || Opt.BaseURL == "" {
		flag.Parse()
	}
}
