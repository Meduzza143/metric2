package handlers

import (
	"bytes"
	"compress/gzip"
	"net/http"

	"github.com/Meduzza143/metric/internal/logger"
)

type (
	responseData struct {
		status int
		size   int
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	//fmt.Printf("buffer sending test:[%v]", string(b))
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func ResponseWritter(w http.ResponseWriter, status int, data []byte, settings RespSettings) {

	l := logger.GetLogger()
	l.Info().Int("status", status).Int("body size", len(data)).Msg("response")

	var answer []byte
	if settings.encoding == "gzip" {
		var buf bytes.Buffer
		gzwriter := gzip.NewWriter(&buf)
		gzwriter.Write(data)
		gzwriter.Close()
		answer = buf.Bytes()
	} else {
		answer = data
	}

	// else {
	// 	answer, _ = io.ReadAll(data)
	// }

	w.WriteHeader(status)
	w.Write(answer)
}
