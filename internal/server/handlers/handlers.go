package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Meduzza143/metric/internal/serializer"
	"github.com/Meduzza143/metric/internal/server/storage"
	"github.com/gorilla/mux"
)

type ExtendedWriter struct {
	http.ResponseWriter
	status         int
	jsonBody       serializer.MetricsJSON
	plainBody      serializer.MetricsPlain
	acceptFormat   string // то в каком виде (структура) передает клиент
	acceptEncoding string // то в каком виде (запакован) принимает клиент
	answerBody     []byte
}

type ExtendedRequester struct {
	req             http.Request
	contentType     string // то в каком виде (структура) принимает клиент
	contentEncoding string // то в каком виде (запакован) передает клиент
	jsonBody        serializer.MetricsJSON
	plainBody       serializer.MetricsPlain
}

/*
			    Сведения о запросах должны содержать URI, метод запроса и время, затраченное на его выполнение.
	    		Сведения об ответах должны содержать код статуса и размер содержимого ответа.
*/

// //http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func UpdateHandle(w http.ResponseWriter, req *http.Request) {

	exWriter, exReq := extend(w, req)
	var answer []byte
	var storage = storage.GetDBHandler()
	defer storage.Close()

	metric, _ := prepareRequest(*exReq)

	if plainValuesCheck(exReq) != nil { //plain text only input check
		exWriter.status = http.StatusBadRequest
		ResponseWritter(*exWriter, []byte("wrong type"))
		return
	}

	exWriter.status = Check(metric.MetricName, metric.MetricType)

	if exWriter.status == http.StatusOK { //ok
		storage.Update(&metric)
		answer = prepareAnswer(*exWriter, metric)
	} else {
		answer = []byte("something went wrong")
	}

	ResponseWritter(*exWriter, answer)

	defer req.Body.Close()
}

func GetMetric(w http.ResponseWriter, req *http.Request) {
	exWriter, exReq := extend(w, req)

	var answer []byte
	var storage = storage.GetDBHandler()
	defer storage.Close()

	metric, err := prepareRequest(*exReq)
	if err != nil {
		exWriter.status = http.StatusBadRequest
		answer = []byte("wrong request")
		ResponseWritter(*exWriter, answer)
		return
	}

	exWriter.status = Check(metric.MetricName, metric.MetricType)

	if exWriter.status != http.StatusOK {
		answer = []byte("wrong metric type")
		ResponseWritter(*exWriter, answer)
		return
	}

	if storage.IsExist(metric.MetricName, metric.MetricType) {
		val := storage.GetOne(metric.MetricName)
		answer = prepareAnswer(*exWriter, *val)
	} else {
		exWriter.status = http.StatusNotFound
		answer = []byte("metric not found")
	}

	ResponseWritter(*exWriter, answer)
	defer req.Body.Close()
}

func GetAll(w http.ResponseWriter, req *http.Request) {
	exWriter, _ := extend(w, req)

	var storage = storage.GetDBHandler()
	defer storage.Close()

	allValues := storage.GetAllValues()

	body := ""
	for k, v := range allValues.GetAllValues() {
		switch v.MetricType {
		case "gauge":
			body += fmt.Sprintf("[%v] %v = %v \n", v.MetricType, k, v.GaugeValue)
		case "counter":
			body += fmt.Sprintf("[%v] %v = %v \n", v.MetricType, k, v.CounterValue)
		}
	}
	ResponseWritter(*exWriter, []byte(body))
	defer req.Body.Close()
}

func PingDB(w http.ResponseWriter, req *http.Request) {
	//При успешной проверке хендлер должен вернуть HTTP-статус 200 OK, при неуспешной — 500 Internal Server Error
	exWriter, _ := extend(w, req)
	if storage.PingPSQL() {
		exWriter.status = http.StatusOK
	} else {
		exWriter.status = http.StatusInternalServerError
	}

	ResponseWritter(*exWriter, []byte(""))
}

func extend(w http.ResponseWriter, req *http.Request) (exWriter *ExtendedWriter, exReq *ExtendedRequester) {
	exReq = &ExtendedRequester{}
	exWriter = &ExtendedWriter{}

	exReq.req = *req
	exWriter.ResponseWriter = w
	exWriter.status = http.StatusOK

	if strings.Contains(req.Header.Get("Content-Encoding"), "gzip") { // what client sending
		exReq.contentEncoding = "gzip"
	}
	if strings.Contains(req.Header.Get("Content-Type"), "application/json") { //what client sending
		exReq.contentType = "json"
	}
	if strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") { //what client can understand
		exWriter.acceptEncoding = "gzip"
	}
	if strings.Contains(req.Header.Get("Accept"), "application/json") { //what client can accept
		exWriter.acceptFormat = "json"
	}
	return
}

func prepareAnswer(exWriter ExtendedWriter, val storage.MemStruct) (answer []byte) {

	switch exWriter.acceptFormat {
	case "json":
		{
			answer = serializer.Serialize(&exWriter.jsonBody, val)
			exWriter.Header().Set("Content-Type", "application/json")
		}
	default:
		{
			answer = serializer.Serialize(&exWriter.plainBody, val)
			exWriter.Header().Set("Content-Type", "text/html")
		}
	}
	return
}

func prepareRequest(exReq ExtendedRequester) (metric storage.MemStruct, err error) {
	if exReq.contentType == "json" {
		metric, err = exReq.jsonBody.Deserialize(&exReq.req)
	} else {
		metric, err = exReq.plainBody.Deserialize(&exReq.req)
	}
	return
}

func plainValuesCheck(exReq *ExtendedRequester) error {
	if exReq.contentType == "json" {
		return nil
	}

	vars := mux.Vars(&exReq.req)
	switch vars["type"] {
	case "gauge":
		{
			_, err := strconv.ParseFloat(vars["value"], 64)
			if err != nil {
				return err
			}
		}
	case "counter":
		{
			_, err := strconv.ParseInt(vars["value"], 10, 64)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
