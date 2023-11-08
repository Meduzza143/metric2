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

type requestCounter int64

var rc *requestCounter = nil

func getCounter() *requestCounter {
	if rc == nil {
		rc = new(requestCounter)
		*rc = 0
	}
	return rc
}

func (rc *requestCounter) incr() {
	*rc += 1
}

func sendData(url, value, name, valueType string) {
	c := getCounter()
	c.incr()
	fmt.Printf("%v ---send number[%v] start---%v\n", strings.Repeat("*", 50), *c, strings.Repeat("*", 50))

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
	for i, v := range request.Header {
		l.Debug().Strs(i, v).Msg("agent set header")
	}

	client := &http.Client{}
	res, err := client.Do(request)

	if err != nil {
		l.Error().Err(err).Msg("got request error")

	} else {
		for i, v := range res.Header {
			l.Debug().Strs(i, v).Msg("agent got header")
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
	fmt.Printf("%v ---send number[%v] end---%v\n", strings.Repeat("#", 50), *c, strings.Repeat("#", 50))

}
