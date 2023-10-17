package server

import (
	"flag"
	"os"
)

type Address struct {
	Listen *string
}

func GetConfig() (a Address) {
	a.initConfig()
	return
}

func (c *Address) initConfig() {

	adr, ok := os.LookupEnv("ADDRESS")
	if !ok {
		flagAdrPtr := flag.String("a", "localhost:8080", "endpont address:port")
		flag.Parse()
		adr = *flagAdrPtr
	}
	c.Listen = &adr
}
