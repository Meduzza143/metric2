package storage

import (
	"net/http"

	"github.com/Meduzza143/metric/internal/logger"
	serverConfig "github.com/Meduzza143/metric/internal/server/settings"
)

type DBHandler interface {
	Update(metric *MemStruct)
	GetOne(name string) *MemStruct
	GetAllValues() MemStorage
	CheckName(name string) int
	checkType(metricType string) int
	Check(name, metrictype string) int
	IsExist(name, metrictype string) bool
}

func GetDBHandler() (db DBHandler) {
	conf := serverConfig.GetConfig()
	if conf.DBType == "PSQL" {
		db = GetPSQLConn()
		return
	} else {
		db = GetMemStorage()
	}
	return
}

func (PSQLData *PostgresSQL) Update(metric *MemStruct) {

}

func (storage MemStorage) Update(metric *MemStruct) {
	metric.CounterValue += storage[metric.MetricName].CounterValue
	storage[metric.MetricName] = *metric

	l := logger.GetLogger()
	l.Info().Str("metric", metric.MetricName).Str("type", metric.MetricType).Int64("counter", metric.CounterValue).Float64("gauge", metric.GaugeValue).Msg("set mem")
}

func (PSQLData *PostgresSQL) GetOne(val string) *MemStruct {
	//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	return &MemStruct{} //mock
}

func (storage MemStorage) GetOne(val string) *MemStruct {
	metric := storage[val]
	return &metric
}

func (PSQLData *PostgresSQL) GetAllValues() MemStorage {
	//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	return make(MemStorage) //mock
}

func (storage MemStorage) GetAllValues() MemStorage {
	return storage
}

// **************************************************************************************************************************************************************************************
func (storage MemStorage) CheckName(name string) (status int) {
	//при успешном приёме возвращать http.StatusOK.
	status = http.StatusOK
	//При попытке передать запрос без имени метрики возвращать http.StatusNotFound.
	if name == "" {
		status = http.StatusNotFound
	}
	return
}

func (PSQLData *PostgresSQL) CheckName(name string) (status int) {
	//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	//TODO::implement
	return
}

func (storage MemStorage) checkType(metricType string) (status int) {
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

func (PSQLData *PostgresSQL) checkType(metricType string) (status int) {
	//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	//TODO::implement
	return
}

func (storage MemStorage) Check(name, metricType string) (status int) {
	status = storage.checkType(metricType)
	if status != http.StatusOK {
		return
	}
	status = storage.CheckName(name)
	return
}

func (PSQLData *PostgresSQL) Check(name, metricType string) (status int) {
	//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	//TODO::implement
	return
}

func (storage MemStorage) IsExist(name, metricType string) bool {
	return storage[name].MetricType == metricType
}

func (PSQLData *PostgresSQL) IsExist(name, metricType string) bool {
	//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	//TODO::implement
	return true
}
