package serializer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Meduzza143/metric/internal/server/storage"
	"github.com/gorilla/mux"
)

type Serializer interface {
	Serialize(storage.MemStruct) []byte
	Deserialize(*http.Request) (storage.MemStruct, error)
}

func Deserialize(s Serializer, req *http.Request) (storage.MemStruct, error) {
	return s.Deserialize(req)
}

func Serialize(s Serializer, mem storage.MemStruct) []byte {
	return s.Serialize(mem)
}

type MetricsJSON struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type MetricsPlain http.Request //Переопределим тушку реквеста чтобы использовать интерфейс

func (*MetricsPlain) Deserialize(req *http.Request) (metric storage.MemStruct, er error) {
	vars := mux.Vars(req)
	metric = storage.MemStruct{}
	metric.MetricType = vars["type"]
	metric.MetricName = vars["name"]
	er = nil

	switch vars["type"] { //gauge : float64, counter: int64
	case "gauge":
		{
			value, _ := strconv.ParseFloat(vars["value"], 64)
			metric.GaugeValue = value
		}
	case "counter":
		{
			value, _ := strconv.ParseInt(vars["value"], 10, 64)
			metric.CounterValue = value
		}
	}
	defer req.Body.Close()
	return
}

func (*MetricsJSON) Deserialize(req *http.Request) (metric storage.MemStruct, er error) {
	var mj MetricsJSON
	er = nil
	body, err := io.ReadAll(req.Body)
	if err != nil {
		er = err
		return
	}
	err = json.Unmarshal(body, &mj)
	if err != nil { //if errror == nil
		er = err
		return
	}
	metric = storage.MemStruct{}

	if mj.ID != "" {
		metric.MetricName = mj.ID
	}
	if mj.MType != "" {
		metric.MetricType = mj.MType
	}
	if mj.Value != nil {
		metric.GaugeValue = *mj.Value
	}
	if mj.Delta != nil {
		metric.CounterValue = *mj.Delta
	}

	defer req.Body.Close()
	return
}

func (*MetricsJSON) Serialize(metric storage.MemStruct) (data []byte) {
	var mj = MetricsJSON{
		MType: metric.MetricType,
		ID:    metric.MetricName,
	}
	switch metric.MetricType {
	case "counter":
		{
			mj.Delta = &metric.CounterValue
		}
	case "gauge":
		{
			mj.Value = &metric.GaugeValue
		}
	}
	data, _ = json.Marshal(mj)
	return
}

func (*MetricsPlain) Serialize(metric storage.MemStruct) (data []byte) {
	switch metric.MetricType {
	case "gauge":
		data = []byte(fmt.Sprintf("%v", metric.GaugeValue))
	case "counter":
		data = []byte(fmt.Sprintf("%v", metric.CounterValue))
	}
	return
}
