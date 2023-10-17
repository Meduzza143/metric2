package main

import (
	"fmt"
	"net/http"

	server "github.com/Meduzza143/metric/internal/server"
	config "github.com/Meduzza143/metric/internal/server/settings"
)

func main() {
	fmt.Println("main test message ... server")
	conf := config.GetConfig()
	r := server.Router()

	fmt.Printf("starting server... at %v \n", conf.Listen)

	err := http.ListenAndServe(conf.Listen, r)
	if err != nil {
		panic(err)
	}
}
