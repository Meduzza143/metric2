package agent

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"

	config "github.com/Meduzza143/metric/internal/agent/config"
	"github.com/Meduzza143/metric/internal/logger"
)

// type reader interface {
// 	ReadAll() ([]byte, error)
// 	Read([]byte) (int, error)
// }

// type readerGzip struct {
// 	rz *gzip.Reader
// 	r  *io.Reader
// }

// func (rg *readerGzip) ReadAll() ([]byte, error) {
// 	return io.ReadAll(rg.rz)
// }

// func (rg *readerGzip) Read(data []byte) (int, error) {
// 	return rg.rz.Read(data)
// }

// func newReader(rw io.Reader) *readerGzip {
// 	gr, _ := gzip.NewReader(rw)
// 	return &readerGzip{gr, &rw}
// }

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

	request, _ = http.NewRequest("POST", finalURL, bytes.NewBuffer(mockData))
	request.Header.Set("Content-Type", "text/plain")

	if cfg.Gzip {
		request.Header.Set("Content-Encoding", "gzip")
	}

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		l.Error().Err(err).Msg("agent")
	}

	var answer []byte
	if cfg.Gzip {
		gzreader, _ := gzip.NewReader(res.Body)
		answer, _ = io.ReadAll(gzreader)
	} else {
		answer, _ = io.ReadAll(res.Body)
	}

	l.Info().Str("answer body", string(answer)).Msg("agent")
	defer res.Body.Close()

	//************************************************************************************************
	// var newReader = newReader(res.Body)
	// test, _ := io.ReadAll(newReader)
	// fmt.Printf("TEST BODY [%v]\n\n", string(test))

	//cannot use r (variable of type *gzip.Reader) as reader value in argument to ReadAll: *gzip.Reader does not implement reader (missing method ReadAll)
	//r, _ := gzip.NewReader(res.Body)
	//ReadAll(r)
	//************************************************************************************************
}
