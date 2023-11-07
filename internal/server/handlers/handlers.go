package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Meduzza143/metric/internal/server/storage"
	"github.com/gorilla/mux"
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
	metric, _ := prepareRequest(w, req)
	//if err == nil {
	status = metric.Check()
	if status == http.StatusOK { //ok
		memStorage.SetValue(&metric)
		answer = prepareAnswer(metric)
	} else {
		answer = []byte("something went wrong")
	}
	// } else {
	// 	status = http.StatusBadRequest
	// }

	ResponseWritter(w, status, answer, respSet)

	defer req.Body.Close()
}

func GetMetric(w http.ResponseWriter, req *http.Request) {
	respSet.Init(req)
	metric, err := prepareRequest(w, req)
	if err == nil {
		status = metric.Check()
		if status == http.StatusOK {
			if metric.IsExist() {
				val := memStorage.GetValue(metric)
				answer = prepareAnswer(val)
			} else {
				status = http.StatusNotFound
			}
		} else {
			answer = []byte("something went wrong")
		}
	} else {
		status = http.StatusBadRequest
	}

	ResponseWritter(w, status, answer, respSet)
	defer req.Body.Close()
}

func GetAll(w http.ResponseWriter, req *http.Request) {
	respSet.Init(req)

	body := ""
	for k, v := range memStorage.GetAllValues() {
		switch v.MetricType {
		case "gauge":
			body += fmt.Sprintf("[%v] %v = %v \n", v.MetricType, k, v.GaugeValue)
		case "counter":
			body += fmt.Sprintf("[%v] %v = %v \n", v.MetricType, k, v.CounterValue)
		}
	}
	ResponseWritter(w, http.StatusOK, []byte(body), respSet)
	defer req.Body.Close()
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

func prepareAnswer(val storage.MemStruct) (answer []byte) {

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
	return
}

func prepareRequest(w http.ResponseWriter, req *http.Request) (metric storage.MemStruct, err error) {
	if respSet.contentType == "json" {
		metric, err = jsonBody.Deserialize(req)
		//w.Header().Set("content-type", "application/json")
	} else {
		er := plainValuesCheck(req)
		if er != nil {
			return
		}
		metric, err = plainBody.Deserialize(req)
		//w.Header().Set("content-type", "text/plain")
	}
	return
}

func plainValuesCheck(req *http.Request) (er error) {
	er = nil
	vars := mux.Vars(req)
	switch vars["type"] {
	case "gauge":
		{
			_, err := strconv.ParseFloat(vars["value"], 64)
			if err != nil {
				er = err
			}
		}
	case "counter":
		{
			_, err := strconv.ParseInt(vars["value"], 10, 64)
			if err != nil {
				er = err
			}
		}
	}
	return
}
