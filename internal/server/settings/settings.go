package server

import (
	"flag"
	"os"
)

type Config struct {
	Listen string
}

func GetConfig() (a *Config) {
	a.initConfig()
	return
}

func (c *Config) initConfig() {

	adr, ok := os.LookupEnv("ADDRESS")
	if !ok {
		flagAdrPtr := flag.String("a", "localhost:8080", "endpont address:port")
		flag.Parse()
		adr = *flagAdrPtr
	}
	c.Listen = adr
}
