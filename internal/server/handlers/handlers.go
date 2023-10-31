package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Meduzza143/metric/internal/logger"
	"github.com/Meduzza143/metric/internal/server/storage"
	"github.com/gorilla/mux"
)

/*
			    Сведения о запросах должны содержать URI, метод запроса и время, затраченное на его выполнение.
	    		Сведения об ответах должны содержать код статуса и размер содержимого ответа.
*/
func LogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		l := logger.GetLogger()
		l.Info().Str("URI", req.URL.Path).Str("Method", req.Method).Str("Remote address", req.RemoteAddr).Msg("request")
		reqStart := time.Now()

		respdata := responseData{
			status: 0,
			size:   0,
		}
		loggingWriter := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   &respdata,
		}

		//next(w, req)
		next(&loggingWriter, req) //какого фига это вообще работает ??????????????????????????????????????????????

		l.Info().Int("status", respdata.status).Int("size", respdata.size).Msg("response")

		reqDuration := time.Now().Sub(reqStart)
		l.Info().Dur("request running time", reqDuration).Msg("request")
	})
}

// //http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func UpdateHandle(w http.ResponseWriter, req *http.Request) {
	var status int
	var metric storage.MemStruct
	var reqIsJson = isJson(*req)
	var answer []byte

	if reqIsJson {
		var jsonBody MetricsJson
		metric, status = jsonBody.Deserialize(req)
		w.Header().Set("content-type", "application/json")
	} else {
		var plainBody MetricsPlain
		metric, status = plainBody.Deserialize(req)
		w.Header().Set("content-type", "text/plain")
	}

	if status == http.StatusOK { //ok
		memStorage := storage.GetInstance()
		memStorage.SetValue(metric.MetricName, metric)

		if reqIsJson {
			var jsonBody MetricsJson
			answer = jsonBody.Serialize(metric)
		} else {
			var plainBody MetricsPlain
			answer = plainBody.Serialize(metric)
		}
	}

	ResponseWritter(w, status, answer) //deserialize answer
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
