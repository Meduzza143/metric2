package main

import (
	"fmt"
	"net/http"

	"github.com/Meduzza143/metric/internal/logger"
	server "github.com/Meduzza143/metric/internal/server"
	"github.com/Meduzza143/metric/internal/server/controllers"
	config "github.com/Meduzza143/metric/internal/server/settings"
)

func main() {
	fmt.Println("getting config...")
	conf := config.GetConfig()
	r := server.Router()
	l := logger.GetLogger()

	s := controllers.GetSaveLoader() //try load data
	if conf.Restore {
		s.LoadAll()
		l.Info().Msg("data has been restored")
	}
	go s.Run()
	defer s.Stop()

	l.Info().Str("Address", conf.Address).Dur("Store interval", conf.StoreInterval).Str("savefile path", conf.StoragePath).
		Bool("restore", conf.Restore).Msg("server starting")

	defer l.Info().Msg("server shut down")

	err := http.ListenAndServe(conf.Address, r)
	if err != nil {
		panic(err)
	}

}
