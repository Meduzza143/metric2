package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Meduzza143/metric/internal/server/storage"
)

type RespSettings struct {
	contentEncoding string // то в каком виде (запакован) передает клиент
	acceptEncoding  string // то в каком виде (запакован) принимает клиент
	acceptFormat    string // то в каком виде (структура) передает клиент
	contentType     string // то в каком виде (структура) принимает клиент
}

var answer []byte
var jsonBody MetricsJson
var plainBody MetricsPlain
var respSet = RespSettings{}
var status int
var metric storage.MemStruct
var memStorage = storage.GetInstance()

/*
			    Сведения о запросах должны содержать URI, метод запроса и время, затраченное на его выполнение.
	    		Сведения об ответах должны содержать код статуса и размер содержимого ответа.
*/

// //http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func UpdateHandle(w http.ResponseWriter, req *http.Request) {

	respSet.Init(req)
	metric, status = prepareRequest(w, req)

	if status == http.StatusOK { //ok
		memStorage.SetValue(metric.MetricName, metric)
		answer, status = prepareAnswer(metric)
	} else {
		answer = []byte("something went wrong")
	}

	ResponseWritter(w, status, answer, respSet)
}

func GetMetric(w http.ResponseWriter, req *http.Request) {
	respSet.Init(req)
	metric, status = prepareRequest(w, req)
	if status == http.StatusOK {
		val := memStorage.GetValue(metric.MetricName)
		answer, status = prepareAnswer(val)
	} else {
		answer = []byte("something went wrong")
	}
	ResponseWritter(w, status, answer, respSet)
}

func GetAll(w http.ResponseWriter, req *http.Request) {
	respSet.Init(req)
	metric, status = prepareRequest(w, req)

	body := ""
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

func (r *RespSettings) Init(req *http.Request) {
	if strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
		r.acceptEncoding = "gzip"
	}
	if strings.Contains(req.Header.Get("Content-Encoding"), "gzip") {
		r.contentEncoding = "gzip"
	}
	if strings.Contains(req.Header.Get("Content-Type"), "application/json") {
		r.contentType = "json"
	}
	if strings.Contains(req.Header.Get("Accept"), "application/json") {
		r.acceptFormat = "json"
	}
}

func prepareAnswer(val storage.MemStruct) (answer []byte, status int) {
	if val.MetricType == metric.MetricType {
		switch respSet.acceptFormat {
		case "json":
			{
				answer = Serialize(&jsonBody, val)
			}
		default:
			{
				answer = Serialize(&plainBody, val)
			}
		}
		status = http.StatusOK
	} else {
		status = http.StatusNotFound
	}
	return
}

func prepareRequest(w http.ResponseWriter, req *http.Request) (metric storage.MemStruct, status int) {
	if respSet.contentType == "json" {
		metric, status = jsonBody.Deserialize(req)
		//w.Header().Set("content-type", "application/json")
	} else {
		metric, status = plainBody.Deserialize(req)
		//w.Header().Set("content-type", "text/plain")
	}
	return
}
