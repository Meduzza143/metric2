package main

import (
	"fmt"
	"net/http"

	"github.com/Meduzza143/metric/internal/logger"
	server "github.com/Meduzza143/metric/internal/server"
	config "github.com/Meduzza143/metric/internal/server/settings"
)

func main() {
	fmt.Println("getting config...")
	conf := config.GetConfig()
	r := server.Router()
	l := logger.GetLogger()

	l.Info().Str("Address", conf.Listen).Msg("server starting")
	defer l.Info().Msg("server shut down")

	err := http.ListenAndServe(conf.Listen, r)
	if err != nil {
		panic(err)
	}

}
