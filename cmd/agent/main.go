package main

import (
	"fmt"
	"time"

	"github.com/Meduzza143/metric/internal/agent"
	config "github.com/Meduzza143/metric/internal/agent/config"
)

func main() {
	fmt.Println("main teset msg ... agent")
	conf := config.GetConfig()
	data := agent.NewStorage()

	reportTicker := time.NewTicker(conf.ReportInterval)
	pollTicker := time.NewTicker(conf.PollInterval)

	fmt.Printf("agent settings:\n address[%v]\n poll interval[%v]\n report interval[%v]", conf.Address, conf.PollInterval, conf.ReportInterval)

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
