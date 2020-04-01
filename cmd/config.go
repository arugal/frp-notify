package main

import (
	"flag"
)

type Config struct {
	Addr          string
	GotifyAddress string
	GotifyToken   string
}

func (c *Config) addFlags() {
	flag.StringVar(&c.Addr, "addr", c.Addr, "server address")
	flag.StringVar(&c.GotifyAddress, "gotify-addr", c.GotifyAddress, "gotify address")
	flag.StringVar(&c.GotifyToken, "gotify-token", c.GotifyToken, "gotify token")
}
