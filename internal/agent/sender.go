package agent

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"

	config "github.com/Meduzza143/metric/internal/agent/config"
)

// type reader interface {
// 	ReadAll() ([]byte, error)
// 	Read([]byte) (int, error)
// }

// type readerPlain struct {
// 	r *io.Reader
// }
// type readerGzip struct {
// 	r *gzip.Reader
// }

// func (rp *readerPlain) ReadAll() ([]byte, error) {
// 	return io.ReadAll(*rp.r)
// }
// func (rg *readerGzip) ReadAll() ([]byte, error) {
// 	return io.ReadAll(rg.r)
// }
// func (rp *readerPlain) Read() ([]byte, error) {
// 	return rp.Read()
// }
// func (rg *readerGzip) Read() ([]byte, error) {
// 	return rg.Read()
// }
// func ReadAll(ri reader) ([]byte, error) {
// 	return io.ReadAll(ri)
// }

func (storage MemStorage) Send(url string) {
	for k, v := range storage {
		sendData(url, v.value, k, v.metricType)
	}
}

func sendData(url, value, name, valueType string) {
	cfg := config.GetConfig()
	var request *http.Request
	finalURL := fmt.Sprintf("%s/update/%s/%s/%s", url, valueType, name, value)

	mockData := []byte(`
	{
		"test": "data",
	}
	`)

	request, err := http.NewRequest("POST", finalURL, bytes.NewBuffer(mockData))
	request.Header.Set("Content-Type", "text/plain")

	if cfg.Gzip {
		request.Header.Set("Content-Encoding", "gzip")
	}

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		fmt.Printf("error [%v]", err)
	} else {
		fmt.Printf("response status [%v]\n", res.Status)
		fmt.Printf("response headers [%v]\n", res.Header)
		gzreader, err := gzip.NewReader(res.Body)
		if err != nil {
			fmt.Println(err)
		}

		output, err2 := io.ReadAll(gzreader)
		if err2 != nil {
			fmt.Println(err2)
		}
		fmt.Printf("body [%v]\n\n", string(output))
	}
	defer res.Body.Close()

	//************************************************************************************************
	//var readerInterface = new(reader)

	//cannot use r (variable of type *gzip.Reader) as reader value in argument to ReadAll: *gzip.Reader does not implement reader (missing method ReadAll)
	// r, _ := gzip.NewReader(res.Body)
	// ReadAll(r)
	//************************************************************************************************
}

// func PackMiddleware(next *bytes.Reader) *bytes.Reader {
// 	var reader bytes.Reader
// 	return &reader
// }
