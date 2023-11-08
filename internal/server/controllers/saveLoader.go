package controllers

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/Meduzza143/metric/internal/logger"
	server "github.com/Meduzza143/metric/internal/server/settings"
	"github.com/Meduzza143/metric/internal/server/storage"
)

type SaveLoader struct {
	path        string
	interval    time.Duration
	keepRunning bool
	file        *os.File
	encoder     *json.Encoder
	decoder     *json.Decoder
}

// type jsonData struct {
// 	ID    string  `json:"id"`              // имя метрики
// 	MType string  `json:"type"`            // параметр, принимающий значение gauge или counter
// 	Delta int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
// 	Value float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
// }

var saveLoader *SaveLoader = nil

func (s *SaveLoader) openWrite() (*os.File, error) {
	return os.OpenFile(s.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_SYNC, 0666)
}
func (s *SaveLoader) openRead() (*os.File, error) {
	return os.OpenFile(s.path, os.O_RDONLY|os.O_SYNC, 0666)
}

func GetSaveLoader() *SaveLoader {
	if saveLoader == nil { //не  больше одного инстанса в одни руки
		conf := server.GetConfig()
		saveLoader = new(SaveLoader)
		saveLoader.interval = conf.StoreInterval
		saveLoader.path = conf.StoragePath

		err := os.MkdirAll(filepath.Dir(saveLoader.path), 0666)
		if err != nil {
			l := logger.GetLogger()
			l.Err(err).Msg("server can't make dir")
		}
	}
	return saveLoader
}

func (s *SaveLoader) Run() {
	s.keepRunning = true
	if s.interval > 0 {
		for s.keepRunning {
			s.SaveAll()
			time.Sleep(s.interval)
		}
		s.file.Close()
	} else if s.interval == 0 {
		//sync input
	}

}

func (s *SaveLoader) Stop() {
	s.keepRunning = false
	s.SaveAll()
	l := logger.GetLogger()
	l.Info().Str("Stopping", "save loader").Msg("server")
}

func (s *SaveLoader) LoadAll() {
	mem := storage.GetInstance()
	l := logger.GetLogger()
	file, err := saveLoader.openRead()
	if err == nil {
		saveLoader.file = file
		saveLoader.decoder = json.NewDecoder(file)
		saveLoader.decoder.Decode(&mem)
		saveLoader.file.Close()
		l.Info().Any("memmory restored", &mem).Msg("server")
	} else {
		l.Info().Err(err).Msg("server can't load data ... initializing db")
		mem.MemInit()
	}

}

func (s *SaveLoader) SaveAll() {
	mem := storage.GetInstance()
	l := logger.GetLogger()

	if len(mem) > 0 {
		file, err := saveLoader.openWrite()
		if err == nil {
			saveLoader.file = file
			saveLoader.encoder = json.NewEncoder(file)
			saveLoader.encoder.Encode(mem)
			saveLoader.file.Close()
		} else {
			l.Info().Err(err).Msg("server")
		}

		l.Info().Str("saved", "data").Msg("server")
	}
}
