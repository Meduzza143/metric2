package main

import (
	"flag"
	"net/http"
	"os"

	server "github.com/Meduzza143/metric/internal/server"
	//config "github.com/Meduzza143/metric/internal/server/settings"
)

func main() {
	//conf := config.GetConfig()
	r := server.Router()

	//fmt.Printf("starting server... at %v \n", conf.Listen)

	//************************************************************************************************
	adr, ok := os.LookupEnv("ADDRESS")
	if !ok {
		flagAdrPtr := flag.String("a", "localhost:8080", "endpont address:port")
		flag.Parse()
		adr = *flagAdrPtr
	}
	//************************************************************************************************

	//err := http.ListenAndServe(conf.Listen, r)
	err := http.ListenAndServe(adr, r)
	if err != nil {
		panic(err)
	}
}
