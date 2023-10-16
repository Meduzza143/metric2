package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Meduzza143/metric/internal/server/storage"
	"github.com/gorilla/mux"
)

// //http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func UpdateHandle(w http.ResponseWriter, req *http.Request) {
	//memStorage := storage.GetInstance()
	memStorage := storage.GetInstance()
	w.Header().Set("content-type", "text/plain")
	vars := mux.Vars(req)

	if vars["name"] == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch vars["type"] {
	case "gauge":
		_, err := strconv.ParseFloat(vars["value"], 64) //оставим проверку на тип
		if err == nil {
			memStorage.SetValue(vars["name"], vars["type"], vars["value"])
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	case "counter":
		_, err := strconv.ParseInt(vars["value"], 0, 64)
		if err == nil { // new value
			thisValue := memStorage.GetValue(vars["name"])
			if (thisValue == storage.MemStruct{}) { //new value
				memStorage.SetValue(vars["name"], vars["type"], vars["value"])
			} else { //increase counter
				currValue, _ := strconv.ParseInt(thisValue.Value, 0, 64)
				currValue += 1
				memStorage.SetValue(vars["name"], vars["type"], strconv.FormatInt(currValue, 10))
			}
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest) //wrong value type
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func GetMetric(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "text/plain")
	memStorage := storage.GetInstance()
	vars := mux.Vars(req)
	val := memStorage.GetValue(vars["name"])
	if val.MetricType == vars["type"] {
		switch val.MetricType {
		case "gauge", "counter":
			w.Write([]byte(fmt.Sprint(val.Value)))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func GetAll(w http.ResponseWriter, req *http.Request) {
	memStorage := storage.GetInstance()
	w.Header().Set("content-type", "text/plain")
	body := "Current values: \n"
	fmt.Println(memStorage.GetAllValues())
	for k, v := range memStorage.GetAllValues() {
		// fmt.Printf("k[%v], v[%v]\n", k, v)
		// fmt.Printf("metric type: %v\n", v.MetricType)
		switch v.MetricType {
		case "gauge":
			body += fmt.Sprintf("%v = %v \n", k, v.Value)
		case "counter":
			body += fmt.Sprintf("%v = %v \n", k, v.Value)
		}
	}
	w.Write([]byte(body))
}
