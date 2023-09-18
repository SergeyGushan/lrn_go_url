package config

import "flag"

type Options struct {
	A string
	B string
}

func (o Options) String() string {
	return o.A + ", " + o.B
}

func (o *Options) Set(s string) error {

	if s == "a" {
		o.A = s
	}

	if s == "b" {
		o.B = s
	}

	return nil
}

func Flags() Options {
	options := Options{}
	flag.StringVar(&options.A, "a", "localhost:8888", "server address")
	flag.StringVar(&options.B, "b", "http://localhost:8000", "redirect address")
	flag.Parse()

	return options
}
