package main

import (
	"time"

	"github.com/Meduzza143/metric/internal/agent"
	config "github.com/Meduzza143/metric/internal/agent/config"
)

func main() {
	conf := config.GetConfig()
	data := agent.NewStorage()

	reportTicker := time.NewTicker(conf.ReportInterval)
	pollTicker := time.NewTicker(conf.PollInterval)

	for {
		select {
		case <-pollTicker.C:
			{
				data.Poll()
				//fmt.Println("poll triggered")
			}
		case <-reportTicker.C:
			{
				data.Send(conf.Address)
				//fmt.Println("report triggered")
			}
		}
	}
}
