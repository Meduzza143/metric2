package logger

import (
	"os"
	"sync"

	zerolog "github.com/rs/zerolog"
)

var zlog *zerolog.Logger = nil
var zlogOnce sync.Once

func GetLogger() *zerolog.Logger {
	zlogOnce.Do(func() {
		var logg = zerolog.New(os.Stdout)
		zlog = &logg
		setConf()
	})

	return zlog
}

func setConf() {
	zlog.Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Logger()
}
