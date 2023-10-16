package agent

import (
	"bytes"
	"fmt"
	"net/http"
)

func (storage MemStorage) Send(url string) {
	for k, v := range storage {
		sendData(url, v.value, k, v.metricType)
	}
}

func sendData(url, value, name, valueType string) {
	finalURL := fmt.Sprintf("%s/update/%s/%s/%s", url, valueType, name, value)
	r := bytes.NewReader([]byte("test"))
	resp, err := http.Post(finalURL, "text/plain", r)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp)
		defer resp.Body.Close()
	}
}
