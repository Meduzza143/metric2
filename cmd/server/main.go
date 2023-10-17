package main

import (
	"flag"
	"fmt"
	"net/http"

	server "github.com/Meduzza143/metric/internal/server"
)

func main() {

	//conf := config.GetConfig()
	r := server.Router()

	//fmt.Printf("starting server... at %v \n", *conf.Listen)

	// var par string
	// flag.StringVar(&par, "a", "localhost:8080", "endpont address:port")
	var ptr = flag.String("a", "localhost:8080", "endpont address:port")
	flag.Parse()

	fmt.Printf("starting server... at %v \n", *ptr)

	err := http.ListenAndServe(*ptr, r)
	if err != nil {
		panic(err)
	}
}
