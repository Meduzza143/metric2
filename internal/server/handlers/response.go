package handlers

import (
	"net/http"

	"github.com/Meduzza143/metric/internal/logger"
	"github.com/Meduzza143/metric/internal/zipper"
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
	if settings.acceptEncoding == "gzip" {
		answer = zipper.GzipBytes(data)
		w.Header().Set("Content-Encoding", "gzip")
	} else {
		answer = data
	}
	w.Header().Set("This server is ", "MINE")
	l.Debug().Str("answer body", string(answer)).Msg("response")

	w.WriteHeader(status)
	w.Write(answer)
}
