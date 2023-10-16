package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"os"
)

type Options struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

var Opt = Options{}

func SetOptions() {
	if err := env.Parse(&Opt); err != nil {
		log.Fatal(err)
	}

	if Opt.ServerAddress == "" {
		if addr := os.Getenv("SERVER_ADDRESS"); addr != "" {
			Opt.ServerAddress = addr
		} else {
			flag.StringVar(&Opt.ServerAddress, "a", "localhost:8080", "server address")
		}
	}

	if Opt.BaseURL == "" {
		if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
			Opt.BaseURL = baseURL
		} else {
			flag.StringVar(&Opt.BaseURL, "b", "http://localhost:8080", "base url")
		}
	}

	flag.Parse()
}
