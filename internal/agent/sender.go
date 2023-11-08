package agent

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

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
		l.Info().Str("sending data", string(zippeddata)).Msg("agent sending body")
		request, _ = http.NewRequest("POST", finalURL, bytes.NewBuffer(zippeddata))
		request.Header.Set("Content-Encoding", "gzip")
		//	request.Header.Set("Accept-Encoding", "gzip")
	} else {
		request, _ = http.NewRequest("POST", finalURL, bytes.NewBuffer(mockData))
		l.Info().Str("sending data", string(mockData)).Msg("agent sending body")
		request.Header.Set("Content-Encoding", "text/plain")
		//	request.Header.Set("Accept-Encoding", "text/plain")
		//request.Header.Set("Accept-Encoding", "identity")
	}
	request.Header.Set("Accept-Encoding", "gzip")

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		l.Error().Err(err).Msg("got request error")
	}

	var answer []byte
	answer, err = io.ReadAll(res.Body)
	if err != nil {
		l.Err(err).Msg("answer read error")
	}
	if strings.Contains(res.Header.Get("Content-Encoding"), "gzip") {
		l.Debug().Str("answer zipped body", string(answer)).Msg("agent zipped response")
		answer = zipper.UnGzipBytes(answer)
	}

	l.Debug().Str("answer body", string(answer)).Msg("agent")
	defer res.Body.Close()
}
