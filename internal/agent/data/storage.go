package data

import (
	"fmt"
	"strconv"
)

type DataStruct struct {
	MetricType   string  `json:"MetricType"`
	GaugeValue   float64 `json:"GaugeValue,omitempty"`
	CounterValue int64   `json:"CounterValue,omitempty"`
	MetricName   string  `json:"id"`
}

type DataStorage map[string]DataStruct

var storage DataStorage = nil

func GetInstance() DataStorage {
	if storage == nil {
		storage = make(DataStorage)
	}
	return storage
}

func (m DataStorage) UpdateMetric(name string, metricType string, value interface{}) {
	newStruct := new(DataStruct)
	newStruct.MetricName = name
	newStruct.MetricType = metricType
	switch metricType {
	case "counter":
		{
			newStruct.CounterValue = value.(int64)
		}
	case "gauge":
		{
			switch value.(type) {
			case float64:
				{
					newStruct.GaugeValue = value.(float64)
				}
			default:
				{
					newStruct.GaugeValue, _ = strconv.ParseFloat(fmt.Sprintf("%v", value), 64)
				}
			}
		}
	}
	m[name] = *newStruct
}
