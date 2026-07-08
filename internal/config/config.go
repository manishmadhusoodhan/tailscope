package config

import (
	"flag"
)

type Config struct {
	Namespace string

	Deployment string

	BufferSize int
}

func Load() Config {

	cfg := Config{}

	flag.StringVar(
		&cfg.Namespace,
		"namespace",
		"default",
		"Kubernetes namespace",
	)

	flag.StringVar(
		&cfg.Deployment,
		"deployment",
		"",
		"Deployment name",
	)

	flag.IntVar(
		&cfg.BufferSize,
		"buffer",
		10000,
		"Log buffer size",
	)

	flag.Parse()

	return cfg
}
