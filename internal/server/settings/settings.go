package serverConfig

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
	DBType        string
	PSQLConn      string
	DBName        string
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
	flagStoreT := flag.Int("i", 300, "store interval: seconds")
	flagFilePath := flag.String("f", "./tmp/metrics-db.json", "file storage path")
	flagRestore := flag.Bool("r", true, "load data file on restart")
	flagPSQLConn := flag.String("d", "", "psql connection")

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

	restore, ok := os.LookupEnv("RESTORE")
	if ok {
		c.Restore, _ = strconv.ParseBool(restore)
	} else {
		c.Restore = *flagRestore
	}

	/*
			При отсутствии переменной окружения DATABASE_DSN или флага командной строки -d или при их пустых значениях вернитесь последовательно к:

		    хранению метрик в файле при наличии соответствующей переменной окружения или флага командной строки;
		    хранению метрик в памяти.
	*/
	c.DBName = "metric"

	PSQLConn, ok := os.LookupEnv("DATABASE_DSN")
	if (ok) && (PSQLConn != "") {
		c.PSQLConn = PSQLConn
		c.DBType = "PSQL"
	} else {
		if *flagPSQLConn != "" {
			c.PSQLConn = *flagPSQLConn
			c.DBType = "PSQL"
		} else { //going mem storage here
			c.DBType = "mem"
		}
	}

}
