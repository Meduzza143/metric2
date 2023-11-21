package storage

import (
	"sync"
)

type MemStruct struct {
	MetricType   string  `json:"MetricType"`
	GaugeValue   float64 `json:"GaugeValue,omitempty"`
	CounterValue int64   `json:"CounterValue,omitempty"`
	MetricName   string  `json:"id"`
}

type MemStorage map[string]MemStruct

var storage MemStorage = nil

func GetMemStorage() MemStorage {
	var storageOnce sync.Once
	storageOnce.Do(func() {
		storage = make(MemStorage)
	})
	return storage
}

// func (metric MemStruct) CheckName() (status int) {
// 	//при успешном приёме возвращать http.StatusOK.
// 	status = http.StatusOK
// 	//При попытке передать запрос без имени метрики возвращать http.StatusNotFound.
// 	if metric.MetricName == "" {
// 		status = http.StatusNotFound
// 	}
// 	return
// }

// func (metric MemStruct) checkType() (status int) {
// 	//при успешном приёме возвращать http.StatusOK.
// 	//При попытке передать запрос с некорректным типом метрики или значением возвращать http.StatusBadRequest
// 	status = http.StatusOK
// 	switch metric.MetricType {
// 	case "gauge", "counter":
// 	default:
// 		status = http.StatusBadRequest
// 	}
// 	// currItem := storage[metric.MetricName]
// 	// if currItem.MetricType != metric.MetricType {
// 	// 	status = http.StatusBadRequest
// 	// }
// 	return
// }

// func (metric MemStruct) Check() (status int) {
// 	status = metric.checkType()
// 	if status == http.StatusOK {
// 		status = metric.CheckName()
// 	}
// 	return
// }

// func (metric MemStruct) IsExist() bool {
// 	currItem := storage[metric.MetricName]
// 	return currItem.MetricType == metric.MetricType
// }
