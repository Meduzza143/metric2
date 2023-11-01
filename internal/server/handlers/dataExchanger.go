package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Meduzza143/metric/internal/server/storage"
	"github.com/gorilla/mux"
)

type MetricsJson struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type MetricsPlain http.Request //Переопределим тушку реквеста чтобы использовать интерфейс

type DataExchanger interface {
	Deserialize() (storage.MemStruct, int)
	Serialize(storage.MemStruct) []byte
	SerializeAll() []byte
}

func (*MetricsPlain) Deserialize(req *http.Request) (metric storage.MemStruct, headerStatus int) {
	vars := mux.Vars(req)
	headerStatus = http.StatusNotFound

	if vars["name"] == "" {
		return
	}

	switch vars["type"] { //gauge : float64, counter: int64
	case "gauge":
		value, err := strconv.ParseFloat(vars["value"], 64) //оставим проверку на тип
		if err == nil {
			metric = storage.MemStruct{
				MetricType: vars["type"],
				GaugeValue: value,
				MetricName: vars["name"],
			}
			headerStatus = http.StatusOK
		} else {
			headerStatus = http.StatusBadRequest
		}
	case "counter":
		value, err := strconv.ParseInt(vars["value"], 10, 64) //оставим проверку на тип
		if err == nil {
			metric = storage.MemStruct{
				MetricType:   vars["type"],
				CounterValue: value,
				MetricName:   vars["name"],
			}
			headerStatus = http.StatusOK
		} else {
			headerStatus = http.StatusBadRequest
		}
	default:
		headerStatus = http.StatusBadRequest
	}
	defer req.Body.Close()
	return metric, headerStatus
}

func (*MetricsJson) Deserialize(req *http.Request) (metric storage.MemStruct, headerStatus int) {
	var mj MetricsJson
	headerStatus = http.StatusBadRequest

	body, err := io.ReadAll(req.Body)
	if err == nil {
		if json.Unmarshal(body, &mj) == nil { //if errror == nil
			metric = storage.MemStruct{}
			switch mj.MType {
			case "gauge":
				{
					metric.MetricName = mj.ID
					metric.MetricType = mj.MType
					metric.GaugeValue = *mj.Value
					headerStatus = http.StatusOK
				}
			case "counter":
				{
					metric.MetricName = mj.ID
					metric.MetricType = mj.MType
					metric.CounterValue = *mj.Delta
					headerStatus = http.StatusOK
				}
			}
		}
	}
	defer req.Body.Close()
	return
}

func (*MetricsJson) Serialize(metric storage.MemStruct) (data []byte) {
	var mj = MetricsJson{
		MType: metric.MetricType,
		ID:    metric.MetricName,
		Delta: &metric.CounterValue,
		Value: &metric.GaugeValue,
	}
	data, _ = json.Marshal(mj)
	return
}

func (*MetricsPlain) Serialize(metric storage.MemStruct) (data []byte) {
	switch metric.MetricType {
	case "gauge":
		data = []byte(fmt.Sprintf("%v = %v \n", metric.MetricName, metric.GaugeValue))
	case "counter":
		data = []byte(fmt.Sprintf("%v = %v \n", metric.MetricName, metric.CounterValue))
	}
	return
}
