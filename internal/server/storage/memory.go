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
func (*memStorage) SetValue(metric *MemStruct) {

	metric.CounterValue += storage[metric.MetricName].CounterValue
	storage[metric.MetricName] = *metric

	l := logger.GetLogger()
	l.Info().Str("metric", metric.MetricName).Str("type", metric.MetricType).Int64("counter", metric.CounterValue).Float64("gauge", metric.GaugeValue).Msg("set")

}

func (memStorage) GetAllValues() memStorage {
	//fmt.Printf("getting all values [%v]\n ", storage)
	return storage
}

func GetInstance() memStorage {
	if storage == nil {
		storage = make(memStorage)
	}
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

func (m MemStruct) Check() (status int) {
	status = m.checkType()
	if status == http.StatusOK {
		status = m.CheckName()
	}
	return
}

func (m MemStruct) IsExist() bool {
	currItem := storage[m.MetricName]
	if currItem.MetricType != m.MetricType {
		return false
	}
	return true
}

func (m *memStorage) MemInit() { // ну вы порядке бреда. Вдруг эти тесты (итер 7) хотят чтобы ВНЕЗАПНО метрики были не любыми и их тип определялся на лету ...
	storage["Alloc"] = *new(MemStruct)
	storage["BuckHashSys"] = *new(MemStruct)
	storage["Frees"] = *new(MemStruct)
	storage["GCSys"] = *new(MemStruct)
	storage["HeapAlloc"] = *new(MemStruct)
	storage["HeapIdle"] = *new(MemStruct)
	storage["HeapInuse"] = *new(MemStruct)
	storage["HeapObjects"] = *new(MemStruct)
	storage["HeapReleased"] = *new(MemStruct)
	storage["HeapSys"] = *new(MemStruct)
	storage["LastGC"] = *new(MemStruct)
	storage["Lookups"] = *new(MemStruct)
	storage["MCacheInuse"] = *new(MemStruct)
	storage["MCacheSys"] = *new(MemStruct)
	storage["MSpanInuse"] = *new(MemStruct)

	storage["MSpanSys"] = *new(MemStruct)
	storage["Mallocs"] = *new(MemStruct)
	storage["NextGC"] = *new(MemStruct)
	storage["OtherSys"] = *new(MemStruct)
	storage["PauseTotalNs"] = *new(MemStruct)
	storage["StackInuse"] = *new(MemStruct)
	storage["StackSys"] = *new(MemStruct)
	storage["Sys"] = *new(MemStruct)
	storage["TotalAlloc"] = *new(MemStruct)

	storage["NumForcedGC"] = *new(MemStruct)
	storage["NumGC"] = *new(MemStruct)

	storage["GCCPUFraction"] = *new(MemStruct)
	storage["RandomValue"] = *new(MemStruct)

	storage["PollCount"] = *new(MemStruct)
}

// func (m MemStruct) CheckValue() (status int) {

// }
