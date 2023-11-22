package storage

import (
	"database/sql"

	"github.com/Meduzza143/metric/internal/logger"
	config "github.com/Meduzza143/metric/internal/server/settings"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //pq driver
)

type PostgresSQL struct {
	db *sqlx.DB
}

// Для хранения значений gauge рекомендуется использовать тип double precision
type DBMetric struct {
	MetricType   string          `db:"MetricType"`
	GaugeValue   sql.NullFloat64 `db:"GaugeValue"`
	CounterValue sql.NullInt64   `db:"CounterValue"`
	MetricName   string          `db:"MetricName"`
}

func connect(conn string) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", conn+"?sslmode=disable&client_encoding=UTF8")
}

func GetPSQLConn() *PostgresSQL {
	l := logger.GetLogger()
	s := config.GetConfig()
	con := new(PostgresSQL)
	db, err := connect(s.PSQLConn + "/" + s.DBName)
	if err != nil {
		l.Err(err).Msg("PSQL connection failed")
		return con
	} else {
		l.Info().Msg("PSQL connection opened")
		con.db = db
	}
	return con
}

func InitPSQLStorage(conn string) (err error) {
	l := logger.GetLogger()
	PSQLdb, err := connect(conn)
	if err != nil {
		l.Err(err).Msg("PSQL connection failed")
		return
	}
	_, err = PSQLdb.Exec(`CREATE DATABASE metric`)
	if err != nil {
		l.Err(err).Msg("DB creation failed")
		//return
	}
	PSQLdb.Close()

	conn += "/metric"
	PSQLdb, err = connect(conn)
	defer PSQLdb.Close()

	if err != nil {
		l.Err(err).Msg("PSQL DB connect failed")
		return
	}

	q := `CREATE TABLE IF NOT EXISTS metric (
		MetricType text,
		GaugeValue double precision,
		CounterValue numeric,
		MetricName text,
		UNIQUE (MetricName)
	);`

	_, err = PSQLdb.Exec(q)
	if err != nil {
		l.Err(err).Msg("PSQL table creation failed")
		return
	}
	return
}

func PingPSQL() (isAvailable bool) {
	s := config.GetConfig()
	conn, err := connect(s.PSQLConn + "/" + s.DBName)
	conn.Close()
	return err == nil
}
