package handlers

import (
	"net/http"
	"strconv"
	"strings"
)

type MemStruct struct {
	metricType   string
	gaugeValue   float64
	counterValue int64
	// gauge   float64 //значение должно замещать
	// counter int64   //значение должно суммироваться
}

//  func updater (s string, k string, t interface{}){
// 	switch t.(type) {
// 		case string

//  }

// func (g float64) updater(s string, k string) {
// 	currValue, err := strconv.ParseFloat(s, 64)
// 	if err == nil {
// 		MemStorage[k].value = currValue
// 	}
// }

/*
	This variable m is a map of string keys to int values:

	var m map[string]int

*/

// var MemStorage map[string]MemStruct //Name - key
var MemStorage = make(map[string]MemStruct)

//m := make(map[string]float64)

// //http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func UpdateHandle(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "text/plain")
	reqPar := strings.Split(strings.TrimPrefix(req.RequestURI, "/update/"), "/")

	if len(reqPar) < 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch reqPar[0] {
	case "gauge":
		currValue, err := strconv.ParseFloat(reqPar[2], 64)
		if err == nil {
			MemStorage[reqPar[1]] = MemStruct{reqPar[0], currValue, 0}
			w.WriteHeader(http.StatusOK)
			//w.Write([]byte("ok"))
		} else {
			w.WriteHeader(http.StatusBadRequest)
			//w.Write([]byte("fail"))
		}
	case "counter":
		currValue, err := strconv.ParseInt(reqPar[2], 0, 64)
		if err == nil {
			thisValue := MemStorage[reqPar[1]].counterValue
			MemStorage[reqPar[1]] = MemStruct{reqPar[0], 0, currValue + thisValue}
			w.WriteHeader(http.StatusOK)
			//w.Write([]byte("ok"))
		} else {
			w.WriteHeader(http.StatusBadRequest)
			//w.Write([]byte("fail"))
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		//w.Write([]byte("fail"))
	}
}
