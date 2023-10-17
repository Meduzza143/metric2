package server

import (
	"flag"
	"os"
)

type Config struct {
	Address string
}

func (c *Config) GetConfig() *Config {
	c.initConfig()
	return c
}

func (c *Config) initConfig() {

	adr, ok := os.LookupEnv("ADDRESS")
	if !ok {
		flagAdrPtr := flag.String("a", "localhost:8080", "endpont address:port")
		flag.Parse()
		adr = *flagAdrPtr
	}
	c.Address = adr
}
