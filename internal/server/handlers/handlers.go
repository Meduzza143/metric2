package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Meduzza143/metric/internal/logger"
	"github.com/Meduzza143/metric/internal/server/storage"
	"github.com/gorilla/mux"
)

// //////// testing mux.Use
// func LogMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
// 		l := logger.GetLogger()
// 		l.Info().Msg("test log middleware")

// 		// compare the return-value to the authMW
// 		next.ServeHTTP(w, req)
// 	})
// }

/*
			    Сведения о запросах должны содержать URI, метод запроса и время, затраченное на его выполнение.
	    		Сведения об ответах должны содержать код статуса и размер содержимого ответа.
*/
func LogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		l := logger.GetLogger()
		l.Info().Str("URI", req.URL.Path).Str("Method", req.Method).Str("Remote address", req.RemoteAddr).Msg("request")
		reqStart := time.Now()

		next(w, req)

		reqDuration := time.Now().Sub(reqStart)
		l.Info().Dur("request running time", reqDuration).Msg("request")
	})
}

// //http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func UpdateHandle(w http.ResponseWriter, req *http.Request) {
	memStorage := storage.GetInstance()
	w.Header().Set("content-type", "text/plain")
	vars := mux.Vars(req)
	headerStatus := http.StatusNotFound

	if vars["name"] == "" {
		return
	}

	switch vars["type"] { //gauge : float64, counter: int64
	case "gauge":
		_, err := strconv.ParseFloat(vars["value"], 64) //оставим проверку на тип
		if err == nil {
			memStorage.SetValue(vars["name"], vars["type"], vars["value"])
			headerStatus = http.StatusOK
		} else {
			headerStatus = http.StatusBadRequest
		}
	case "counter":
		val, err := strconv.ParseInt(vars["value"], 0, 64)
		if err == nil { // new value
			thisValue := memStorage.GetValue(vars["name"])
			if (thisValue == storage.MemStruct{}) { //new value
				memStorage.SetValue(vars["name"], vars["type"], vars["value"])
			} else { //increase counter
				currValue, _ := strconv.ParseInt(thisValue.Value, 0, 64)
				currValue += val
				memStorage.SetValue(vars["name"], vars["type"], strconv.FormatInt(currValue, 10))
			}
			headerStatus = http.StatusOK
		} else {
			headerStatus = http.StatusBadRequest
		}
	default:
		headerStatus = http.StatusBadRequest
	}
	ResponseWritter(w, headerStatus, "")
}

func GetMetric(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "text/plain")
	memStorage := storage.GetInstance()
	vars := mux.Vars(req)
	val := memStorage.GetValue(vars["name"])
	headerStatus := http.StatusNotFound
	if val.MetricType == vars["type"] {
		switch val.MetricType {
		case "gauge", "counter":
			ResponseWritter(w, http.StatusOK, fmt.Sprint(val.Value))
		default:
			ResponseWritter(w, headerStatus, "")
		}
	} else {
		ResponseWritter(w, headerStatus, "")
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
			body += fmt.Sprintf("%v = %v \n", k, v.Value)
		case "counter":
			body += fmt.Sprintf("%v = %v \n", k, v.Value)
		}
	}
	ResponseWritter(w, http.StatusOK, body)
	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte(body))
}
