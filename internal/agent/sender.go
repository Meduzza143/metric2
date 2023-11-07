package agent

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	config "github.com/Meduzza143/metric/internal/agent/config"
	"github.com/Meduzza143/metric/internal/logger"
	"github.com/Meduzza143/metric/internal/zipper"
)

func (storage MemStorage) Send(url string) {
	for k, v := range storage {
		sendData(url, v.value, k, v.metricType)
	}
}

func sendData(url, value, name, valueType string) {
	cfg := config.GetConfig()
	var request *http.Request
	l := logger.GetLogger()

	finalURL := fmt.Sprintf("%s/update/%s/%s/%s", url, valueType, name, value)

	l.Info().Str("sending", finalURL).Msg("agent")

	mockData := []byte(`
	{
		"test": "data",
	}
	`)

	//cfg.Gzip = true //TEST!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

	if cfg.Gzip {
		zippeddata := zipper.GzipBytes(mockData) //zipper.GzipBytes(request.Body)
		request, _ = http.NewRequest("POST", finalURL, bytes.NewBuffer(zippeddata))
		request.Header.Set("Content-Encoding", "gzip")
		request.Header.Set("Accept-Encoding", "gzip")
	} else {
		request, _ = http.NewRequest("POST", finalURL, bytes.NewBuffer(mockData))
		request.Header.Set("Content-Type", "text/plain")
		request.Header.Set("Accept-Encoding", "identity")
	}

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		l.Error().Err(err).Msg("agent")
	}

	var answer []byte
	answer, _ = io.ReadAll(res.Body)
	if cfg.Gzip {
		answer = zipper.UnGzipBytes(answer)
	}

	l.Debug().Str("answer body", string(answer)).Msg("agent")
	defer res.Body.Close()
}
