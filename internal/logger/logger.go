package logger

import (
	"os"

	zerolog "github.com/rs/zerolog"
)

var zlog *zerolog.Logger = nil

func GetLogger() *zerolog.Logger {
	if zlog == nil {
		var logg = zerolog.New(os.Stdout)
		zlog = &logg
		setConf()
	}
	return zlog
}

func setConf() {
	zlog.Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Logger()
}
