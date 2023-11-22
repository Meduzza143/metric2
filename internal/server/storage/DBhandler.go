package storage

import (
	"fmt"

	"github.com/Meduzza143/metric/internal/logger"
	serverConfig "github.com/Meduzza143/metric/internal/server/settings"
)

type DBHandler interface {
	Update(metric *MemStruct)
	GetOne(name string) *MemStruct
	GetAllValues() MemStorage
	IsExist(name, metrictype string) bool
	Close()
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
	l := logger.GetLogger()
	// q := `
	// INSERT INTO metric
	// SET MetricType = $1, GaugeValue = $2, CounterValue = $3
	// WHERE MetricName = $4;
	// `
	q := `INSERT INTO metric (metricType, gaugeValue, counterValue, metricName) VALUES ($1, $2, $3, $4) ON CONFLICT (metricName) DO UPDATE
		SET MetricType = $1, GaugeValue = $2, CounterValue = $3`
	_, err := PSQLData.db.Exec(q, metric.MetricType, metric.GaugeValue, metric.CounterValue, metric.MetricName)
	if err != nil {
		l.Err(err).Str("name", metric.MetricName).Msg("metric update failed")
	} else {
		l.Info().Str("result", metric.MetricName).Msg("metric updated successfully")
	}
}

func (storage MemStorage) Update(metric *MemStruct) {
	metric.CounterValue += storage[metric.MetricName].CounterValue
	storage[metric.MetricName] = *metric

	l := logger.GetLogger()
	l.Info().Str("metric", metric.MetricName).Str("type", metric.MetricType).Int64("counter", metric.CounterValue).Float64("gauge", metric.GaugeValue).Msg("set mem")
}

func (PSQLData *PostgresSQL) GetOne(name string) *MemStruct {
	metric := new(MemStruct)
	PSQLData.db.Select(metric, "SELECT * FROM metric WHERE MetricName=$1", name)
	return metric
}

func (storage MemStorage) GetOne(name string) *MemStruct {
	metric := storage[name]
	return &metric
}

func (PSQLData *PostgresSQL) GetAllValues() (answer MemStorage) {
	data := []MemStruct{}
	err := PSQLData.db.Select(data, "SELECT * FROM metric")

	for _, v := range data {
		answer[v.MetricName] = v
	}

	if err != nil {
		fmt.Println(err)
	}
	return
}

func (storage MemStorage) GetAllValues() MemStorage {
	return storage
}

func (storage MemStorage) IsExist(name, metricType string) bool {
	return storage[name].MetricType == metricType
}

func (PSQLData *PostgresSQL) IsExist(name, metricType string) bool {
	row := PSQLData.db.QueryRow("SELECT * FROM metric WHERE MetricName=$1, MetricType=$2", name, metricType)
	err := row.Scan()
	return err == nil
}

func (storage MemStorage) Close() {
}

func (PSQLData *PostgresSQL) Close() {
	//	PSQLData.db.Close()
}
