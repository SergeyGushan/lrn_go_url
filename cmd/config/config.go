package config

import "flag"

type Options struct {
	A string
	B string
}

func Flags() Options {
	options := Options{}
	flag.StringVar(&options.A, "a", "localhost:8888", "server address")
	flag.StringVar(&options.B, "b", "http://localhost:8000", "redirect address")
	flag.Parse()

	return options
}
