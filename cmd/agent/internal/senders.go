package senders

import (
	"bytes"
	"fmt"
	"net/http"
)

var APIURL = "http://127.0.0.1:8080/update"

func SendData(value string, name string, valueType string) int {
	//var ret int = 0
	finalURL := fmt.Sprintf("%s/%s/%s/%s", APIURL, valueType, name, value)
	r := bytes.NewReader([]byte("test"))
	resp, err := http.Post(finalURL, "text/plain", r)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp)
	}
	//ret = resp.StatusCode
	return resp.StatusCode
}
