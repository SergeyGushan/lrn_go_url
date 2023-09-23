package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

type Options struct {
	A string `env:"SERVER_ADDRESS"`
	B string `env:"BASE_URL"`
}

var Opt = Options{}

func SetOptions() {
	err := env.Parse(&Opt)
	if err != nil {
		log.Fatal(err)
	}

	if Opt.A == "" {
		flag.StringVar(&Opt.A, "a", "localhost:8080", "server address")
	}
	if Opt.B == "" {
		flag.StringVar(&Opt.B, "b", "http://localhost:8080", "base url")
	}

	if Opt.A == "" || Opt.B == "" {
		flag.Parse()
	}
}
