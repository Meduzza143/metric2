package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Meduzza143/metric/internal/logger"
	server "github.com/Meduzza143/metric/internal/server"
	"github.com/Meduzza143/metric/internal/server/controllers"
	config "github.com/Meduzza143/metric/internal/server/settings"
	"github.com/xlab/closer"
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
		l.Info().Msg("restoring data ...")
		s.LoadAll()
	}

	ctx, cancel := context.WithCancel(context.Background())
	go s.Run(ctx)
	closer.Bind(func() { //перехватывает системные события на закрытие программы (не обязательно аккуратно) и вызывает функцию
		l.Info().Msg("server shut down")
		cancel()
	})

	err := http.ListenAndServe(conf.Address, r)
	if err != nil {
		panic(err)
	}
}
