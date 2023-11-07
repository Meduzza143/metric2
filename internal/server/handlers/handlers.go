package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Meduzza143/metric/internal/server/storage"
)

type RespSettings struct {
	contentEncoding string
	acceptEncoding  string
	format          string
}

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
	//var reqIsJson = isJson(*req)

	respSet := RespSettings{}
	respSet.Init(req)

	if respSet.format == "json" {
		metric, status = jsonBody.Deserialize(req)
		w.Header().Set("content-type", "application/json")
	} else {
		metric, status = plainBody.Deserialize(req)
		w.Header().Set("content-type", "text/plain")
	}

	if status == http.StatusOK { //ok
		memStorage := storage.GetInstance()
		memStorage.SetValue(metric.MetricName, metric)
		val := memStorage.GetValue(metric.MetricName)
		switch respSet.format {
		case "json":
			{
				answer = jsonBody.Serialize(val)
			}
		default:
			{
				answer = plainBody.Serialize(val)
			}
		}
	}

	ResponseWritter(w, status, answer, respSet)
}

func GetMetric(w http.ResponseWriter, req *http.Request) {

	var status int
	var metric storage.MemStruct
	var answer []byte
	var jsonBody MetricsJson
	var plainBody MetricsPlain

	respSet := RespSettings{}
	respSet.Init(req)

	var reqIsJson = isJson(*req)
	if reqIsJson {
		metric, status = jsonBody.Deserialize(req)
		w.Header().Set("content-type", "application/json")
	} else {
		metric, status = plainBody.Deserialize(req)
		w.Header().Set("content-type", "text/plain")
	}

	//if status == http.StatusOK { //ok
	memStorage := storage.GetInstance()
	val := memStorage.GetValue(metric.MetricName)
	if val.MetricType == metric.MetricType {
		switch respSet.format {
		case "json":
			{
				answer = jsonBody.Serialize(val)
			}
		default:
			{
				answer = plainBody.Serialize(val)
			}
		}
		status = http.StatusOK
	} else {
		status = http.StatusNotFound
	}
	//}

	ResponseWritter(w, status, answer, respSet)
}

func GetAll(w http.ResponseWriter, req *http.Request) {
	respSet := RespSettings{}
	respSet.Init(req)

	memStorage := storage.GetInstance()
	w.Header().Set("content-type", "text/plain")

	body := "Current values: \n"

	for k, v := range memStorage.GetAllValues() {
		switch v.MetricType {
		case "gauge":
			body += fmt.Sprintf("%v = %v \n", k, v.GaugeValue)
		case "counter":
			body += fmt.Sprintf("%v = %v \n", k, v.CounterValue)
		}
	}
	ResponseWritter(w, http.StatusOK, []byte(body), respSet)
}

func isJson(r http.Request) bool {
	if r.Header.Get("Content-Type") == "application/json" {
		return true
	}
	return false
}

func (r *RespSettings) Init(req *http.Request) {
	if strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
		r.acceptEncoding = "gzip"
	}
	if strings.Contains(req.Header.Get("Content-Encoding"), "gzip") {
		r.contentEncoding = "gzip"
	}
}
