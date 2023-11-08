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

	s := controllers.GetSaveLoader()

	l.Info().Str("Address", conf.Address).Dur("Store interval", conf.StoreInterval).Str("savefile path", conf.StoragePath).
		Bool("restore", conf.Restore).Msg("server starting")

	if conf.Restore {
		l.Info().Msg("restorinmg data ...")
		s.LoadAll()
	}
	go s.Run()
	defer s.Stop()

	defer l.Info().Msg("server shut down")

	err := http.ListenAndServe(conf.Address, r)
	if err != nil {
		panic(err)
	}

}
