package storage

import (
	"net/http"
	"sync"

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
var storageOnce sync.Once

func (*memStorage) SetValue(metric *MemStruct) {

	metric.CounterValue += storage[metric.MetricName].CounterValue
	storage[metric.MetricName] = *metric

	l := logger.GetLogger()
	l.Info().Str("metric", metric.MetricName).Str("type", metric.MetricType).Int64("counter", metric.CounterValue).Float64("gauge", metric.GaugeValue).Msg("set")
}

func (memStorage) GetAllValues() memStorage {
	return storage
}

func GetInstance() memStorage {
	storageOnce.Do(func() {
		storage = make(memStorage)
	})
	return storage
}

func (memStorage) GetValue(val MemStruct) (answer MemStruct) {
	answer = storage[val.MetricName]
	return
}

func (metric MemStruct) CheckName() (status int) {
	//при успешном приёме возвращать http.StatusOK.
	status = http.StatusOK
	//При попытке передать запрос без имени метрики возвращать http.StatusNotFound.
	if metric.MetricName == "" {
		status = http.StatusNotFound
	}
	return
}

func (metric MemStruct) checkType() (status int) {
	//при успешном приёме возвращать http.StatusOK.
	//При попытке передать запрос с некорректным типом метрики или значением возвращать http.StatusBadRequest
	status = http.StatusOK
	switch metric.MetricType {
	case "gauge", "counter":
	default:
		status = http.StatusBadRequest
	}
	// currItem := storage[metric.MetricName]
	// if currItem.MetricType != metric.MetricType {
	// 	status = http.StatusBadRequest
	// }
	return
}

func (metric MemStruct) Check() (status int) {
	status = metric.checkType()
	if status == http.StatusOK {
		status = metric.CheckName()
	}
	return
}

func (metric MemStruct) IsExist() bool {
	currItem := storage[metric.MetricName]
	return currItem.MetricType == metric.MetricType
}
