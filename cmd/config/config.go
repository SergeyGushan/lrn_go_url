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

func SetOptions() Options {
	options := Options{}

	err := env.Parse(&options)
	if err != nil {
		log.Fatal(err)
	}

	if options.A == "" {
		flag.StringVar(&options.A, "a", "localhost:8888", "server address")
	}
	if options.B == "" {
		flag.StringVar(&options.B, "b", "http://localhost:8000", "base url")
	}

	if options.A == "" || options.B == "" {
		flag.Parse()
	}

	return options
}
