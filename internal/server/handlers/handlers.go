package handlers

import (
	"fmt"
	"net/http"

	"github.com/Meduzza143/metric/internal/server/storage"
	"github.com/gorilla/mux"
)

/*
			    Сведения о запросах должны содержать URI, метод запроса и время, затраченное на его выполнение.
	    		Сведения об ответах должны содержать код статуса и размер содержимого ответа.
*/

// //http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func UpdateHandle(w http.ResponseWriter, req *http.Request) {
	var status int
	var metric storage.MemStruct
	var answer []byte
	var jsonBody MetricsJson
	var plainBody MetricsPlain

	var reqIsJson = isJson(*req)
	if reqIsJson {
		metric, status = jsonBody.Deserialize(req)
		w.Header().Set("content-type", "application/json")
	} else {
		metric, status = plainBody.Deserialize(req)
		w.Header().Set("content-type", "text/plain")
	}

	if status == http.StatusOK { //ok
		memStorage := storage.GetInstance()
		memStorage.SetValue(metric.MetricName, metric)

		if reqIsJson {
			answer = jsonBody.Serialize(metric)
		} else {
			answer = plainBody.Serialize(metric)
		}
	}

	ResponseWritter(w, status, answer)
}

func GetMetric(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("content-type", "text/plain")
	memStorage := storage.GetInstance()
	vars := mux.Vars(req)
	val := memStorage.GetValue(vars["name"])
	//headerStatus := http.StatusNotFound
	if val.MetricType == vars["type"] {
		switch val.MetricType {
		case "gauge", "counter":
			//ResponseWritter(w, http.StatusOK, fmt.Sprint(val.Value))
			//ResponseWritter(w, http.StatusOK, fmt.Sprint(val.GaugeValue))
		default:
			//ResponseWritter(w, headerStatus, "")
		}
	} else {
		//		ResponseWritter(w, headerStatus, "")
	}
}

func GetAll(w http.ResponseWriter, req *http.Request) {
	memStorage := storage.GetInstance()
	w.Header().Set("content-type", "text/plain")

	body := "Current values: \n"
	//fmt.Println(memStorage.GetAllValues())
	for k, v := range memStorage.GetAllValues() {
		switch v.MetricType {
		case "gauge":
			//body += fmt.Sprintf("%v = %v \n", k, v.Value)
			body += fmt.Sprintf("%v = %v \n", k, v.GaugeValue)
		case "counter":
			//body += fmt.Sprintf("%v = %v \n", k, v.Value)
			body += fmt.Sprintf("%v = %v \n", k, v.CounterValue)
		}
	}
	//	ResponseWritter(w, http.StatusOK, body)
	//
	// w.WriteHeader(http.StatusOK)
	//w.Write([]byte(body))
}

func isJson(r http.Request) bool {
	if r.Header.Get("Content-Type") == "application/json" {
		return true
	}
	return false
}
