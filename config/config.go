package config

import "flag"

var cfg *Config

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
}

func New() *Config {
	cfg = &Config{}
	flag.StringVar(&cfg.Host, "h", "localhost", "NATS server host")
	flag.IntVar(&cfg.Port, "p", 4222, "NATS server port")
	flag.StringVar(&cfg.Username, "u", "", "NATS server username")
	flag.StringVar(&cfg.Password, "P", "", "NATS server password")

	flag.Parse()

	return cfg
}

func GetConfig() *Config {
	return cfg
}
