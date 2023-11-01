package main

import (
	"fmt"
	"time"

	"github.com/Meduzza143/metric/internal/agent"
	config "github.com/Meduzza143/metric/internal/agent/config"
)

func main() {
	fmt.Println("starting agent ...")
	cfg := config.GetConfig()
	data := agent.NewStorage()

	// if conf.CheckConfig() == false {
	// 	fmt.Printf("wrong arguments:")
	// 	fmt.Printf("agent settings:\n address[%v]\n poll interval[%v]\n report interval[%v]", conf.Address, conf.PollInterval, conf.ReportInterval)
	// 	os.Exit(0)
	// }

	fmt.Printf("agent settings:\n address[%v]\n poll interval[%v]\n report interval[%v]", cfg.Address, cfg.PollInterval, cfg.ReportInterval)

	reportTicker := time.NewTicker(cfg.ReportInterval)
	fmt.Println("report ticker has been set")
	pollTicker := time.NewTicker(cfg.PollInterval)
	fmt.Println("poll ticker has been set")

	for {
		select {
		case <-pollTicker.C:
			{
				data.Poll()
				//fmt.Println("poll triggered")
			}
		case <-reportTicker.C:
			{
				data.Send(cfg.Address)
				//fmt.Println("report triggered")
			}
		}
	}
}
