package main

import (
	"fmt"
	"time"

	"github.com/Meduzza143/metric/internal/agent"
	config "github.com/Meduzza143/metric/internal/agent/config"
)

func main() {
	fmt.Println("starting agent ...")
	conf := config.GetConfig()
	data := agent.NewStorage()

	// if conf.CheckConfig() == false {
	// 	fmt.Printf("wrong arguments:")
	// 	fmt.Printf("agent settings:\n address[%v]\n poll interval[%v]\n report interval[%v]", conf.Address, conf.PollInterval, conf.ReportInterval)
	// 	os.Exit(0)
	// }

	fmt.Printf("agent settings:\n address[%v]\n poll interval[%v]\n report interval[%v]", conf.Address, conf.PollInterval, conf.ReportInterval)

	reportTicker := time.NewTicker(conf.ReportInterval)
	fmt.Println("report ticker has been set")
	pollTicker := time.NewTicker(conf.PollInterval)
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
				data.Send(conf.Address)
				//fmt.Println("report triggered")
			}
		}
	}
}
