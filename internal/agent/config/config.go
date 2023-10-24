package agent

import (
	"flag"
	"os"
	"time"
)

type Settings struct {
	Address string
	//port           string
	//
	PollInterval   time.Duration
	ReportInterval time.Duration
}

func GetConfig() (c Settings) {
	c.initConfig()
	return
}

func (c *Settings) initConfig() {

	flagAdrPtr := flag.String("a", "localhost:8080", "endpont address:port")
	flagRepPtr := flag.Int("r", 10, "report interval in seconds")
	flagPolPtr := flag.Int("p", 2, "poll interval in seconds")

	flag.Parse()

	adr, ok := os.LookupEnv("ADDRESS")
	if !ok {
		adr = *flagAdrPtr
	}
	c.Address = "http://" + adr

	repString, ok := os.LookupEnv("REPORT_INTERVAL")
	if !ok {
		c.ReportInterval = time.Duration(*flagRepPtr) * time.Second
	} else {
		c.ReportInterval, _ = time.ParseDuration(repString + "s")
	}

	pollString, ok := os.LookupEnv("POLL_INTERVAL")
	if !ok {
		c.PollInterval = time.Duration(*flagPolPtr) * time.Second
	} else {
		c.PollInterval, _ = time.ParseDuration(pollString + "s")
	}
}

// func (c *Settings) CheckConfig() (pass bool) {
// 	if c.ReportInterval <= 0 {
// 		pass = false
// 	}
// 	//poll interval check
// 	if c.PollInterval <= 0 {
// 		pass = false
// 	}

// 	return
// }
