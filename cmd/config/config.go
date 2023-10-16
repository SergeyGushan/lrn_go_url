package config

import (
	"flag"
	"os"
)

type Options struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
	ServerPort    string `env:"SERVER_PORT"`
}

var Opt = Options{}

func SetOptions() {
	if Opt.ServerAddress == "" {
		if addr := os.Getenv("SERVER_ADDRESS"); addr != "" {
			Opt.ServerAddress = addr
		} else {
			flag.StringVar(&Opt.ServerAddress, "server-address", "localhost", "server address")
		}
	}

	if Opt.ServerPort == "" {
		if port := os.Getenv("SERVER_PORT"); port != "" {
			Opt.ServerPort = port
		} else {
			flag.StringVar(&Opt.ServerPort, "server-port", "8080", "server port")
		}
	}

	if Opt.BaseURL == "" {
		if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
			Opt.BaseURL = baseURL
		} else {
			flag.StringVar(&Opt.BaseURL, "base-url", "http://localhost:8080", "base url")
		}
	}

	flag.Parse()
}
