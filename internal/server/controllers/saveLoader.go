package controllers

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Meduzza143/metric/internal/logger"
	serverConfig "github.com/Meduzza143/metric/internal/server/settings"
	"github.com/Meduzza143/metric/internal/server/storage"
)

type SaveLoader struct {
	path     string
	interval time.Duration
	file     *os.File
	encoder  *json.Encoder
	decoder  *json.Decoder
}

var saveLoader *SaveLoader = nil
var saveLoaderOnce sync.Once

func (s *SaveLoader) openWrite() (*os.File, error) {
	return os.OpenFile(s.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_SYNC, 0666)
}
func (s *SaveLoader) openRead() (*os.File, error) {
	return os.OpenFile(s.path, os.O_RDONLY|os.O_SYNC, 0666)
}

func GetSaveLoader() *SaveLoader {
	saveLoaderOnce.Do(func() {
		conf := serverConfig.GetConfig()
		saveLoader = new(SaveLoader)
		saveLoader.interval = conf.StoreInterval
		saveLoader.path = conf.StoragePath

		err := os.MkdirAll(filepath.Dir(saveLoader.path), 0666)
		if err != nil {
			l := logger.GetLogger()
			l.Err(err).Msg("server can't make dir")
		}
	})
	return saveLoader
}

func (s *SaveLoader) Run(ctx context.Context) {
	if s.interval <= 0 {
		return
	}

	t := time.NewTicker(s.interval)
	for {
		select {
		case <-t.C:
			s.SaveAll()
		case <-ctx.Done():
			s.SaveAll()
			s.file.Close()
			return
		}
	}
}

func (s *SaveLoader) LoadAll() {
	mem := storage.GetMemStorage()
	l := logger.GetLogger()
	file, err := saveLoader.openRead()
	if err != nil {
		l.Info().Err(err).Msg("server can't load data ... new mem struct")
		return
	}
	saveLoader.file = file
	saveLoader.decoder = json.NewDecoder(file)
	saveLoader.decoder.Decode(&mem)
	saveLoader.file.Close()
	l.Info().Any("memmory restored", &mem).Msg("server")
}

func (s *SaveLoader) SaveAll() {
	mem := storage.GetMemStorage()
	l := logger.GetLogger()

	if len(mem) > 0 {
		file, err := saveLoader.openWrite()
		if err != nil {
			l.Info().Err(err).Msg("server")
			return
		}
		saveLoader.file = file
		saveLoader.encoder = json.NewEncoder(file)
		saveLoader.encoder.Encode(mem)
		saveLoader.file.Close()
		l.Info().Str("saved", "data").Msg("server")
	}
}
