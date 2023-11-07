package server

import (
	"flag"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Address       string
	StoreInterval time.Duration
	StoragePath   string
	Restore       bool
}

var c *Config = nil

func GetConfig() Config {
	if c == nil {
		c = new(Config)
		c.initConfig()
	}
	return *c
}

func (c *Config) initConfig() {

	flagAdrPtr := flag.String("a", "localhost:8080", "endpont address:port")
	flagStoreT := flag.Int("i", 3, "store interval: seconds")
	flagFilePath := flag.String("f", "./tmp/metrics-db.json", "file storage path")
	//flagFilePath := flag.String("f", "metrics-db.json", "file storage path")
	flagRestore := flag.Bool("r", true, "load data file on restart")
	flag.Parse()

	adr, ok := os.LookupEnv("ADDRESS")
	if ok {
		c.Address = adr
	} else {
		c.Address = *flagAdrPtr
	}

	store, ok := os.LookupEnv("STORE_INTERVAL")
	if ok {
		c.StoreInterval, _ = time.ParseDuration(store + "s")
	} else {
		c.StoreInterval = time.Duration(*flagStoreT) * time.Second
	}

	storagePath, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if ok {
		c.StoragePath = storagePath
	} else {
		c.StoragePath = *flagFilePath
	}

	restore, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if ok {
		c.Restore, _ = strconv.ParseBool(restore)
	} else {
		c.Restore = *flagRestore
	}

}
