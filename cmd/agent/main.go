package main

import (
	"fmt"
	"time"

	"github.com/Meduzza143/metric/internal/agent"
	config "github.com/Meduzza143/metric/internal/agent/config"
	"github.com/Meduzza143/metric/internal/logger"
)

func main() {
	fmt.Println("starting agent ...")
	cfg := config.GetConfig()
	data := agent.NewStorage()
	l := logger.GetLogger()

	l.Info().Str("server address", cfg.Address).Msg("Agent")
	l.Info().Dur("report interval", cfg.ReportInterval).Dur("poll interval", cfg.PollInterval).Bool("use gzip", cfg.Gzip).Msg("Agent starting")

	reportTicker := time.NewTicker(cfg.ReportInterval)
	pollTicker := time.NewTicker(cfg.PollInterval)

	for {
		select {
		case <-pollTicker.C:
			{
				data.Poll()
			}
		case <-reportTicker.C:
			{
				data.Send(cfg.Address)
			}
		}
	}
}
