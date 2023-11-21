package storage

import (
	"github.com/Meduzza143/metric/internal/logger"
	config "github.com/Meduzza143/metric/internal/server/settings"
	"github.com/jmoiron/sqlx"
)

type PostgresSQL struct {
	db *sqlx.DB
}

// Для хранения значений gauge рекомендуется использовать тип double precision
type DBMetric struct {
	MetricType   string  `db:"MetricType"`
	GaugeValue   float64 `db:"GaugeValue"`
	CounterValue int64   `db:"CounterValue"`
	MetricName   string  `db:"MetricName"`
}

func GetPSQLConn() (PSQLConn *PostgresSQL) {
	l := logger.GetLogger()
	s := config.GetConfig()
	db, err := sqlx.Connect("postgres", s.PSQLConn)
	if err != nil {
		l.Err(err).Msg("PSQL connection failed")
		return
	} else {
		l.Info().Msg("PSQL connection opened")
		PSQLConn.db = db
	}
	return
}

func InitPSQLStorage(conn string) (err error) {
	l := logger.GetLogger()
	conn = "postgres://postgres:postgres@localhost:5433" //psql connection string
	PSQLdb, _ := sqlx.Connect("postgres", conn+"?sslmode=disable&client_encoding=UTF8")
	PSQLdb.Exec(`CREATE DATABASE metric`)
	PSQLdb.Close()

	conn += "/metric/"
	PSQLdb, err = sqlx.Connect("postgres", conn)
	if err != nil {
		l.Err(err).Msg("PSQL connection failed")
		return
	}

	q := `CREATE TABLE metric (
		MetricType text,
		GaugeValue double precision,
		CounterValue numeric,
		MetricName text
	);`

	PSQLdb.Exec(q)
	return
}

func PingPSQL() (isAvailable bool) {
	conf := config.GetConfig()
	_, err := sqlx.Connect("postgres", conf.PSQLConn+"?sslmode=disable&client_encoding=UTF8")
	return err == nil
}
