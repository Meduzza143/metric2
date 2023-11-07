package storage

import (
	"net/http"

	"github.com/Meduzza143/metric/internal/logger"
)

type MemStruct struct {
	MetricType string `json:"MetricType"`
	//Value      string
	GaugeValue   float64 `json:"GaugeValue,omitempty"`
	CounterValue int64   `json:"CounterValue,omitempty"`
	MetricName   string  `json:"id"`
}

type memStorage map[string]MemStruct

var storage memStorage = nil

//TODO: implement interface

// type Storage interface {
// 	SetValue(key, metricType, value string)
// 	GetAllValues() memStorage
// 	GetValue(key string) MemStruct
// }

// func (memStorage) SetValue(MetricName, metricType string, gaugeValue float64, counterValue int64) {
func (memStorage) SetValue(MetricName string, metric MemStruct) {

	metric.CounterValue = storage[MetricName].CounterValue + metric.CounterValue

	storage[MetricName] = metric
	l := logger.GetLogger()
	l.Info().Str("metric", metric.MetricName).Str("type", metric.MetricType).Int64("counter", metric.CounterValue).Float64("gauge", metric.GaugeValue).Msg("set")

}

func (memStorage) GetAllValues() memStorage {
	//fmt.Printf("getting all values [%v]\n ", storage)
	return storage
}

func (memStorage) GetValue(val MemStruct) (answer MemStruct, status int) {

	currItem := storage[val.MetricName]
	if currItem.MetricType == val.MetricType {
		answer = currItem
		status = http.StatusOK
	} else {
		status = http.StatusNotFound
	}
	return
}

func GetInstance() memStorage {
	if storage == nil {
		storage = make(memStorage)
	}
	return storage
}
