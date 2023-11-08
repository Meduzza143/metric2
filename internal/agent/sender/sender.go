package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Meduzza143/metric/internal/agent/config"
	"github.com/Meduzza143/metric/internal/agent/data"
	"github.com/Meduzza143/metric/internal/logger"
	"github.com/Meduzza143/metric/internal/serializer"
	"github.com/Meduzza143/metric/internal/zipper"
)

func Send(url string) {
	storage := data.GetInstance()
	for _, v := range storage {
		sendData(url, v)
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

func sendData(url string, metric data.DataStruct) {
	c := getCounter()
	c.incr()
	fmt.Printf("%v ---send number[%v] start---%v\n", strings.Repeat("*", 50), *c, strings.Repeat("*", 50))

	cfg := config.GetConfig()
	var request *http.Request
	l := logger.GetLogger()

	// var value string
	// if metric.MetricType == "counter" {
	// 	value = strconv.FormatInt(metric.CounterValue, 10)
	// } else { //gauge
	// 	value = strconv.FormatFloat(metric.GaugeValue, 'f', -1, 64)
	// }
	// finalURL := fmt.Sprintf("%s/update/%s/%s/%s", url, metric.MetricType, metric.MetricName, value)

	// l.Info().Str("sending", finalURL).Msg("agent")

	var finalURL string = url + "/update/"

	var mj = serializer.MetricsJson{
		MType: metric.MetricType,
		ID:    metric.MetricName,
		Delta: &metric.CounterValue,
		Value: &metric.GaugeValue,
	}

	data, _ := json.Marshal(mj)

	if cfg.Gzip {
		zippeddata := zipper.GzipBytes(data) //zipper.GzipBytes(request.Body)
		l.Info().Str("sending data", string(zippeddata)).Msg("agent sending body")
		request, _ = http.NewRequest("POST", finalURL, bytes.NewBuffer(zippeddata))
		request.Header.Set("Content-Encoding", "gzip")
	} else {
		request, _ = http.NewRequest("POST", finalURL, bytes.NewBuffer(data))
		l.Info().Str("sending data", string(data)).Msg("agent sending body")
	}
	request.Header.Set("Accept-Encoding", "gzip")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

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
