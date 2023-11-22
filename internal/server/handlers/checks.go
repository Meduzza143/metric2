package handlers

import "net/http"

func CheckName(name string) (status int) {
	//при успешном приёме возвращать http.StatusOK.
	status = http.StatusOK
	//При попытке передать запрос без имени метрики возвращать http.StatusNotFound.
	if name == "" {
		status = http.StatusNotFound
	}
	return
}

func checkType(metricType string) (status int) {
	//при успешном приёме возвращать http.StatusOK.
	//При попытке передать запрос с некорректным типом метрики или значением возвращать http.StatusBadRequest
	status = http.StatusOK
	switch metricType {
	case "gauge", "counter":
	default:
		status = http.StatusBadRequest
	}
	return
}

func Check(name, metricType string) (status int) {
	status = checkType(metricType)
	if status != http.StatusOK {
		return
	}
	status = CheckName(name)
	return
}
