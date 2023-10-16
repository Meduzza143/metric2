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
	flagRepPtr := flag.Duration("r", 10*time.Second, "report interval in seconds")
	flagPolPtr := flag.Duration("p", 2*time.Second, "poll interval in seconds")

	flag.Parse()

	adr, ok := os.LookupEnv("ADDRESS")
	if !ok {
		adr = *flagAdrPtr
	}
	c.Address = "http://" + adr

	repString, ok := os.LookupEnv("REPORT_INTERVAL")
	if !ok {
		c.ReportInterval = *flagRepPtr
	} else {
		c.ReportInterval, _ = time.ParseDuration(repString)
	}

	pollString, ok := os.LookupEnv("POLL_INTERVAL")
	if !ok {
		c.PollInterval = *flagPolPtr
	} else {
		c.PollInterval, _ = time.ParseDuration(pollString)
	}
}
