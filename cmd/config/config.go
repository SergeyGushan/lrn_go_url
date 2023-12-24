package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"log"
	"os"
)

type Options struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

var Opt = Options{}

func SetOptions() {
	flag.StringVar(&Opt.ServerAddress, "a", "localhost:8080", "server address")
	flag.StringVar(&Opt.BaseURL, "b", "http://localhost:8080", "base url")
	flag.StringVar(&Opt.FileStoragePath, "f", os.TempDir()+"/short-url-db.json", "base url")
	flag.StringVar(&Opt.DatabaseDSN, "d", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", `localhost`, `go_user`, `go_password`, `go_learn`), "db dsn")
	flag.Parse()

	err := env.Parse(&Opt)

	if err != nil {
		log.Fatal(err)
	}
}
